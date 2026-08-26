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

func ToCaroMatch(m db.CaroMatch) (CaroMatch, error) {
	var board CaroBoard
	err := json.Unmarshal(m.Board, &board)
	if err != nil {
		return CaroMatch{}, err
	}

	var moves []CaroMove
	err = json.Unmarshal(m.Moves, &moves)
	if err != nil {
		return CaroMatch{}, err
	}

	// X luôn đi trước nên lượt suy ra được từ số nước đã đánh, không cần cột riêng
	turnOf := X
	if len(moves)%2 == 1 {
		turnOf = O
	}

	return CaroMatch{
		Id:                  uuid.UUID(m.ID.Bytes),
		IsRated:             m.IsRated,
		XPlayerId:           uuid.UUID(m.XPlayerID.Bytes),
		XPlayerRatingBefore: int(m.XPlayerRatingBefore),
		XPlayerRatingAfter:  shared.Nullable[int]{Value: int(m.XPlayerRatingAfter.Int32), IsNull: !m.XPlayerRatingAfter.Valid},
		OPlayerId:           uuid.UUID(m.OPlayerID.Bytes),
		OPlayerRatingBefore: int(m.OPlayerRatingBefore),
		OPlayerRatingAfter:  shared.Nullable[int]{Value: int(m.OPlayerRatingAfter.Int32), IsNull: !m.OPlayerRatingAfter.Valid},
		Board:               board,
		Moves:               moves,
		TurnOf:              turnOf,
		Status:              CaroStatus(m.Status),
		EndReason:           CaroEndReason(m.EndReason),
		StartedAt:           m.StartedAt.Time,
		EndedAt:             shared.Nullable[time.Time]{Value: m.EndedAt.Time, IsNull: !m.EndedAt.Valid},
	}, nil
}

func (a *AppState) GetCaroMatchById(ctx context.Context, id uuid.UUID) (CaroMatch, error) {
	m, err := a.q.GetCaroMatchById(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return CaroMatch{}, err
	}
	return ToCaroMatch(m)
}

// ProcessCaroMatchEnded lưu match đã kết thúc xuống DB, cập nhật rating 2 bên,
// dọn match khỏi RAM và trả về chính match đó với rating sau trận đã điền.
func (a *AppState) ProcessCaroMatchEnded(ctx context.Context, m CaroMatch) (CaroMatch, error) {
	result := 0
	winnerId := pgtype.UUID{}
	if m.Status == XWon {
		result = 1
		winnerId = pgtype.UUID{Bytes: m.XPlayerId, Valid: true}
	}
	if m.Status == OWon {
		result = -1
		winnerId = pgtype.UUID{Bytes: m.OPlayerId, Valid: true}
	}

	xRatingChange := calculateRatingChange(m.XPlayerRatingBefore, m.OPlayerRatingBefore, result)
	oRatingChange := calculateRatingChange(m.OPlayerRatingBefore, m.XPlayerRatingBefore, -result)
	m.XPlayerRatingAfter = shared.Nullable[int]{Value: m.XPlayerRatingBefore + xRatingChange}
	m.OPlayerRatingAfter = shared.Nullable[int]{Value: m.OPlayerRatingBefore + oRatingChange}

	encodedBoard, err := json.Marshal(m.Board)
	if err != nil {
		return CaroMatch{}, err
	}
	encodedMoves, err := json.Marshal(m.Moves)
	if err != nil {
		return CaroMatch{}, err
	}

	tx, err := a.p.Begin(ctx)
	if err != nil {
		return CaroMatch{}, err
	}
	defer tx.Rollback(ctx)

	txq := a.q.WithTx(tx)

	err = txq.InsertCaroMatch(ctx, db.InsertCaroMatchParams{
		ID:                  pgtype.UUID{Bytes: m.Id, Valid: true},
		IsRated:             m.IsRated,
		XPlayerID:           pgtype.UUID{Bytes: m.XPlayerId, Valid: true},
		XPlayerRatingBefore: int32(m.XPlayerRatingBefore),
		XPlayerRatingAfter:  pgtype.Int4{Int32: int32(m.XPlayerRatingAfter.Value), Valid: !m.XPlayerRatingAfter.IsNull},
		OPlayerID:           pgtype.UUID{Bytes: m.OPlayerId, Valid: true},
		OPlayerRatingBefore: int32(m.OPlayerRatingBefore),
		OPlayerRatingAfter:  pgtype.Int4{Int32: int32(m.OPlayerRatingAfter.Value), Valid: !m.OPlayerRatingAfter.IsNull},
		WinnerID:            winnerId,
		Status:              string(m.Status),
		EndReason:           string(m.EndReason),
		Board:               encodedBoard,
		Moves:               encodedMoves,
		StartedAt:           pgtype.Timestamptz{Time: m.StartedAt, Valid: true},
		EndedAt:             pgtype.Timestamptz{Time: m.EndedAt.Value, Valid: !m.EndedAt.IsNull},
	})
	if err != nil {
		return CaroMatch{}, err
	}

	err = txq.UpdateUserCaroRating(ctx, db.UpdateUserCaroRatingParams{ID: pgtype.UUID{Bytes: m.XPlayerId, Valid: true}, CaroRating: int32(xRatingChange)})
	if err != nil {
		return CaroMatch{}, err
	}

	err = txq.UpdateUserCaroRating(ctx, db.UpdateUserCaroRatingParams{ID: pgtype.UUID{Bytes: m.OPlayerId, Valid: true}, CaroRating: int32(oRatingChange)})
	if err != nil {
		return CaroMatch{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return CaroMatch{}, err
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

	return m, nil
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
