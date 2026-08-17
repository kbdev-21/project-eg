package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const CARO_BOARD_SIZE = 15
const CARO_MAX_MOVE_TIME = 20 * time.Second

type CaroMatch struct {
	Id uuid.UUID `json:"id"`

	IsRated bool `json:"isRated"`

	XPlayerId     pgtype.UUID `json:"xPlayerId"`
	XPlayerRating int         `json:"xPlayerRating"`
	OPlayerId     pgtype.UUID `json:"oPlayerId"`
	OPlayerRating int         `json:"oPlayerRating"`

	Board  CaroBoard  `json:"board"`
	Moves  []CaroMove `json:"moves"`
	TurnOf CaroPiece  `json:"turnOf"`

	Winner    CaroPiece `json:"winner"`
	IsEnded   bool      `json:"isEnded"`
	StartedAt time.Time `json:"startedAt"`
}

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

func NewCaroMatch(id uuid.UUID, isRated bool, xId pgtype.UUID, xRating int, oId pgtype.UUID, oRating int) *CaroMatch {
	return &CaroMatch{
		Id:            id,
		IsRated:       isRated,
		XPlayerId:     xId,
		XPlayerRating: xRating,
		OPlayerId:     oId,
		OPlayerRating: oRating,
		Board:         CaroBoard{},
		Moves:         []CaroMove{},
		TurnOf:        X,
		Winner:        None,
		IsEnded:       false,
		StartedAt:     time.Now(),
	}
}

func (match *CaroMatch) Move(player CaroPiece, x int, y int) error {
	if player != X && player != O {
		return fmt.Errorf("invalid input")
	}
	if match.IsEnded {
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
		match.IsEnded = true
		return nil
	}

	// check win
	isWin, err := checkWinCondition(validMove.Piece, validMove.X, validMove.Y, match.Board)
	if err != nil {
		return err
	}

	if isWin {
		match.Winner = validMove.Piece
		match.IsEnded = true
	}

	return nil
}

func (match *CaroMatch) OutOfTime() {
	ootPlayer := match.TurnOf
	if ootPlayer == X {
		match.Winner = O
	} else {
		match.Winner = X
	}
	match.IsEnded = true
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
