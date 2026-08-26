package domain

import (
	"backend/src/shared"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const CARO_BOARD_SIZE = 15
const CARO_MAX_MOVE_TIME = 20 * time.Second

// CaroMatch dùng chung cho match đang chơi (trong RAM) và match cũ (đọc từ DB).
// Field nullable = chỉ biết sau khi trận kết thúc.
type CaroMatch struct {
	Id uuid.UUID `json:"id"`

	IsRated bool `json:"isRated"`

	XPlayerId           uuid.UUID            `json:"xPlayerId"`
	XPlayerRatingBefore int                  `json:"xPlayerRatingBefore"`
	XPlayerRatingAfter  shared.Nullable[int] `json:"xPlayerRatingAfter"`
	OPlayerId           uuid.UUID            `json:"oPlayerId"`
	OPlayerRatingBefore int                  `json:"oPlayerRatingBefore"`
	OPlayerRatingAfter  shared.Nullable[int] `json:"oPlayerRatingAfter"`

	Board  CaroBoard  `json:"board"`
	Moves  []CaroMove `json:"moves"`
	TurnOf CaroPiece  `json:"turnOf"`

	Status    CaroStatus                 `json:"status"`
	EndReason CaroEndReason              `json:"endReason"`
	StartedAt time.Time                  `json:"startedAt"`
	EndedAt   shared.Nullable[time.Time] `json:"endedAt"`
}

type CaroStatus string

const (
	Playing CaroStatus = "PLAYING"
	XWon    CaroStatus = "X_WON"
	OWon    CaroStatus = "O_WON"
	Draw    CaroStatus = "DRAW"
)

type CaroEndReason string

const (
	NotEnded  CaroEndReason = ""
	FiveInRow CaroEndReason = "FIVE_IN_ROW"
	FullBoard CaroEndReason = "FULL_BOARD"
	OutOfTime CaroEndReason = "OUT_OF_TIME"
)

type CaroBoard [CARO_BOARD_SIZE][CARO_BOARD_SIZE]CaroPiece

type CaroPiece string

const (
	None CaroPiece = ""
	X    CaroPiece = "X"
	O    CaroPiece = "O"
)

type CaroMove struct {
	Piece    CaroPiece `json:"piece"`
	X        int       `json:"x"`
	Y        int       `json:"y"`
	PlayedAt time.Time `json:"playedAt"`
}

func NewCaroMatch(id uuid.UUID, isRated bool, xId uuid.UUID, xRating int, oId uuid.UUID, oRating int) *CaroMatch {
	return &CaroMatch{
		Id:                  id,
		IsRated:             isRated,
		XPlayerId:           xId,
		XPlayerRatingBefore: xRating,
		XPlayerRatingAfter:  shared.Nullable[int]{IsNull: true},
		OPlayerId:           oId,
		OPlayerRatingBefore: oRating,
		OPlayerRatingAfter:  shared.Nullable[int]{IsNull: true},
		Board:               CaroBoard{},
		Moves:               []CaroMove{},
		TurnOf:              X,
		Status:              Playing,
		EndReason:           NotEnded,
		StartedAt:           time.Now(),
		EndedAt:             shared.Nullable[time.Time]{IsNull: true},
	}
}

func (match *CaroMatch) Move(player CaroPiece, x int, y int) error {
	if player != X && player != O {
		return fmt.Errorf("invalid input")
	}
	if match.Status != Playing {
		return fmt.Errorf("match ended")
	}
	if match.TurnOf != player {
		return fmt.Errorf("invalid move")
	}
	if x < 0 || x >= CARO_BOARD_SIZE || y < 0 || y >= CARO_BOARD_SIZE {
		return fmt.Errorf("invalid move")
	}
	if match.Board[y][x] != None {
		return fmt.Errorf("invalid move")
	}

	validMove := CaroMove{
		Piece:    player,
		X:        x,
		Y:        y,
		PlayedAt: time.Now(),
	}

	match.Board[y][x] = validMove.Piece
	match.Moves = append(match.Moves, validMove)
	if match.TurnOf == X {
		match.TurnOf = O
	} else {
		match.TurnOf = X
	}

	// check draw
	if len(match.Moves) == CARO_BOARD_SIZE*CARO_BOARD_SIZE {
		match.Status = Draw
		match.EndReason = FullBoard
		match.EndedAt = shared.Nullable[time.Time]{Value: validMove.PlayedAt}
		return nil
	}

	// check win
	isWin, err := checkWinCondition(validMove.Piece, validMove.X, validMove.Y, match.Board)
	if err != nil {
		return err
	}

	if isWin {
		match.Status = XWon
		if validMove.Piece == O {
			match.Status = OWon
		}
		match.EndReason = FiveInRow
		match.EndedAt = shared.Nullable[time.Time]{Value: validMove.PlayedAt}
	}

	return nil
}

func (match *CaroMatch) OutOfTime() {
	match.Status = OWon
	if match.TurnOf == O {
		match.Status = XWon
	}
	match.EndReason = OutOfTime
	match.EndedAt = shared.Nullable[time.Time]{Value: time.Now()}
}

var caroDirections = [4][2]int{
	{1, 0},  // ngang
	{0, 1},  // dọc
	{1, 1},  // chéo xuống-phải
	{1, -1}, // chéo lên-phải
}

func checkWinCondition(piece CaroPiece, fromX int, fromY int, board CaroBoard) (bool, error) {
	if piece != X && piece != O {
		return false, fmt.Errorf("invalid piece")
	}
	if fromX < 0 || fromX >= CARO_BOARD_SIZE || fromY < 0 || fromY >= CARO_BOARD_SIZE {
		return false, fmt.Errorf("position out of bounds")
	}
	if board[fromY][fromX] != piece {
		return false, fmt.Errorf("piece at position does not match")
	}

	for _, dir := range caroDirections {
		count := 1

		for _, sign := range [2]int{1, -1} {
			dx := dir[0] * sign
			dy := dir[1] * sign

			x := fromX + dx
			y := fromY + dy
			for x >= 0 && x < CARO_BOARD_SIZE && y >= 0 && y < CARO_BOARD_SIZE && board[y][x] == piece {
				count++
				x += dx
				y += dy
			}
		}

		if count >= 5 {
			return true, nil
		}
	}

	return false, nil
}
