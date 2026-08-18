package domain

import (
	"backend/src/db"
	"context"
	"encoding/json"
	"math"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CaroMatchResult struct {
	db.CaroMatch
	FinalBoard CaroBoard  `json:"finalBoard"`
	Moves      []CaroMove `json:"moves"`
}

func ToCaroMatchResult(m db.CaroMatch) CaroMatchResult {
	var finalBoard CaroBoard
	json.Unmarshal(m.FinalBoard, &finalBoard)

	var moves []CaroMove
	json.Unmarshal(m.Moves, &moves)

	return CaroMatchResult{
		CaroMatch:  m,
		FinalBoard: finalBoard,
		Moves:      moves,
	}
}

func (a *AppState) GetCaroMatchResultById(ctx context.Context, id pgtype.UUID) (CaroMatchResult, error) {
	m, err := a.q.GetCaroMatchById(ctx, id)
	if err != nil {
		return CaroMatchResult{}, err
	}
	return ToCaroMatchResult(m), nil
}

func (a *AppState) ProcessCaroMatchEnded(ctx context.Context, m CaroMatch) (CaroMatchResult, error) {
	result := 0
	winnerId := pgtype.UUID{}
	if m.Winner == X {
		result = 1
		winnerId = pgtype.UUID{Bytes: m.XPlayerId.Bytes, Valid: true}
	}
	if m.Winner == O {
		result = -1
		winnerId = pgtype.UUID{Bytes: m.OPlayerId.Bytes, Valid: true}
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

	dbMatchId := pgtype.UUID{Bytes: m.Id, Valid: true}
	err = txq.InsertCaroMatch(ctx, db.InsertCaroMatchParams{
		ID:                  dbMatchId,
		XPlayerID:           m.XPlayerId,
		XPlayerRatingBefore: int32(m.XPlayerRating),
		XPlayerRatingAfter:  int32(m.XPlayerRating + xRatingChange),
		OPlayerID:           m.OPlayerId,
		OPlayerRatingBefore: int32(m.OPlayerRating),
		OPlayerRatingAfter:  int32(m.OPlayerRating + oRatingChange),
		WinnerID:            winnerId,
		FinalBoard:          encodedBoard,
		Moves:               encodedMoves,
	})
	if err != nil {
		return CaroMatchResult{}, err
	}

	err = txq.UpdateUserCaroRating(ctx, db.UpdateUserCaroRatingParams{ID: m.XPlayerId, CaroRating: int32(xRatingChange)})
	if err != nil {
		return CaroMatchResult{}, err
	}

	err = txq.UpdateUserCaroRating(ctx, db.UpdateUserCaroRatingParams{ID: m.OPlayerId, CaroRating: int32(oRatingChange)})
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

	return a.GetCaroMatchResultById(ctx, dbMatchId)
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
