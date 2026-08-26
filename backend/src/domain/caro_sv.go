package domain

import (
	"backend/src/db"
	"backend/src/shared"
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CaroMatchResult struct {
	Id                  uuid.UUID                  `json:"id"`
	XPlayerId           uuid.UUID                  `json:"xPlayerId"`
	XPlayerRatingBefore int                        `json:"xPlayerRatingBefore"`
	XPlayerRatingAfter  int                        `json:"xPlayerRatingAfter"`
	OPlayerId           uuid.UUID                  `json:"oPlayerId"`
	OPlayerRatingBefore int                        `json:"oPlayerRatingBefore"`
	OPlayerRatingAfter  int                        `json:"oPlayerRatingAfter"`
	WinnerId            shared.Nullable[uuid.UUID] `json:"winnerId"`
	FinalBoard          CaroBoard                  `json:"finalBoard"`
	Moves               []CaroMove                 `json:"moves"`
	CreatedAt           time.Time                  `json:"createdAt"`
	UpdatedAt           time.Time                  `json:"updatedAt"`
}

func ToCaroMatchResult(m db.CaroMatch) (CaroMatchResult, error) {
	var finalBoard CaroBoard
	err := json.Unmarshal(m.FinalBoard, &finalBoard)
	if err != nil {
		return CaroMatchResult{}, err
	}

	var moves []CaroMove
	err = json.Unmarshal(m.Moves, &moves)
	if err != nil {
		return CaroMatchResult{}, err
	}

	return CaroMatchResult{
		Id:                  uuid.UUID(m.ID.Bytes),
		XPlayerId:           uuid.UUID(m.XPlayerID.Bytes),
		XPlayerRatingBefore: int(m.XPlayerRatingBefore),
		XPlayerRatingAfter:  int(m.XPlayerRatingAfter),
		OPlayerId:           uuid.UUID(m.OPlayerID.Bytes),
		OPlayerRatingBefore: int(m.OPlayerRatingBefore),
		OPlayerRatingAfter:  int(m.OPlayerRatingAfter),
		WinnerId:            shared.Nullable[uuid.UUID]{Value: uuid.UUID(m.WinnerID.Bytes), IsNull: !m.WinnerID.Valid},
		FinalBoard:          finalBoard,
		Moves:               moves,
		CreatedAt:           m.CreatedAt.Time,
		UpdatedAt:           m.UpdatedAt.Time,
	}, nil
}

func (a *AppState) GetCaroMatchResultById(ctx context.Context, id uuid.UUID) (CaroMatchResult, error) {
	m, err := a.q.GetCaroMatchById(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return CaroMatchResult{}, err
	}
	return ToCaroMatchResult(m)
}

func (a *AppState) ProcessCaroMatchEnded(ctx context.Context, m CaroMatch) (CaroMatchResult, error) {
	result := 0
	winnerId := pgtype.UUID{}
	if m.Winner == X {
		result = 1
		winnerId = pgtype.UUID{Bytes: m.XPlayerId, Valid: true}
	}
	if m.Winner == O {
		result = -1
		winnerId = pgtype.UUID{Bytes: m.OPlayerId, Valid: true}
	}
	xRatingChange := calculateRatingChange(m.XPlayerRating, m.OPlayerRating, result)
	oRatingChange := calculateRatingChange(m.OPlayerRating, m.XPlayerRating, -result)
	encodedBoard, _ := json.Marshal(m.Board)
	encodedMoves, _ := json.Marshal(m.Moves)

	tx, err := a.p.Begin(ctx)
	if err != nil {
		return CaroMatchResult{}, err
	}
	defer tx.Rollback(ctx)

	txq := a.q.WithTx(tx)

	err = txq.InsertCaroMatch(ctx, db.InsertCaroMatchParams{
		ID:                  pgtype.UUID{Bytes: m.Id, Valid: true},
		XPlayerID:           pgtype.UUID{Bytes: m.XPlayerId, Valid: true},
		XPlayerRatingBefore: int32(m.XPlayerRating),
		XPlayerRatingAfter:  int32(m.XPlayerRating + xRatingChange),
		OPlayerID:           pgtype.UUID{Bytes: m.OPlayerId, Valid: true},
		OPlayerRatingBefore: int32(m.OPlayerRating),
		OPlayerRatingAfter:  int32(m.OPlayerRating + oRatingChange),
		WinnerID:            winnerId,
		FinalBoard:          encodedBoard,
		Moves:               encodedMoves,
	})
	if err != nil {
		return CaroMatchResult{}, err
	}

	err = txq.UpdateUserCaroRating(ctx, db.UpdateUserCaroRatingParams{ID: pgtype.UUID{Bytes: m.XPlayerId, Valid: true}, CaroRating: int32(xRatingChange)})
	if err != nil {
		return CaroMatchResult{}, err
	}

	err = txq.UpdateUserCaroRating(ctx, db.UpdateUserCaroRatingParams{ID: pgtype.UUID{Bytes: m.OPlayerId, Valid: true}, CaroRating: int32(oRatingChange)})
	if err != nil {
		return CaroMatchResult{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return CaroMatchResult{}, err
	}

	a.mu.Lock()

	delete(a.caroMatchesMap, m.Id)

	xSes, ok := a.userSessionsMap[m.XPlayerId]
	if ok {
		xSes.CurrentMatchId = uuid.Nil
		xSes.State = Idle
	}
	oSes, ok := a.userSessionsMap[m.OPlayerId]
	if ok {
		oSes.CurrentMatchId = uuid.Nil
		oSes.State = Idle
	}

	a.mu.Unlock()

	return a.GetCaroMatchResultById(ctx, m.Id)
}

// result: -1 lose, 0 draw, 1 win
func calculateRatingChange(rating int, opRating int, result int) int {
	if rating <= 0 {
		rating = 1
	}
	if opRating <= 0 {
		opRating = 1
	}

	const (
		winLoseBase = 50.0
		drawSwing   = 10.0
	)

	diff := math.Log2(float64(opRating) / float64(rating))
	diff = math.Min(math.Max(diff, -1), 1)

	switch {
	case result > 0:
		return int(math.Round(winLoseBase + diff*winLoseBase/2))
	case result < 0:
		return int(math.Round(-winLoseBase + diff*winLoseBase/2))
	default:
		return int(math.Round(diff * drawSwing))
	}
}
