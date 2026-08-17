package router

import (
	"backend/src/config"
	"backend/src/db"
	"backend/src/domain"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientMessage struct {
	Code ClientMessageCode `json:"code"`
	Data json.RawMessage   `json:"data"`
}

type ClientMessageCode string

const (
	Ping           ClientMessageCode = "PING"
	CaroJoinQueue  ClientMessageCode = "CARO:JOIN_QUEUE"
	CaroLeaveQueue ClientMessageCode = "CARO:LEAVE_QUEUE"
	CaroPlayMove   ClientMessageCode = "CARO:PLAY_MOVE"
	CaroAbortMatch ClientMessageCode = "CARO:ABORT_MATCH"
)

type ServerMessage struct {
	Code  ServerMessageCode       `json:"code"`
	State domain.UserSessionState `json:"state"`
	Data  any                     `json:"data"`
}

type ServerMessageCode string

const (
	Ok                      ServerMessageCode = "OK"
	Error                   ServerMessageCode = "ERROR"
	CaroMatchFound          ServerMessageCode = "CARO:MATCH_FOUND"
	CaroNewBoardState       ServerMessageCode = "CARO:NEW_BOARD_STATE"
	CaroMatchEnded          ServerMessageCode = "CARO:MATCH_ENDED"
	CaroMatchEndedOutOfTime ServerMessageCode = "CARO:MATCH_ENDED_OUT_OF_TIME"
)

func InitWsRoute(fib *fiber.App, q *db.Queries, conf config.Config, app *domain.AppState, pool *pgxpool.Pool) {
	fib.Get("/ws", wsAuthMiddleware(q, conf, pool), websocket.New(func(wsc *websocket.Conn) {
		uId := wsc.Locals("currentUser").(domain.User).ID

		// save userSession
		userSession := app.GetOrCreateUserSession(uId)
		if userSession.WsConn != nil {
			userSession.WsConn.Close()
		}
		userSession.WsConn = wsc
		log.Printf("Websocket: User %s connect", uId)

		for {
			_, data, err := wsc.ReadMessage()
			// disconnect
			if err != nil {
				wsc.Close()
				if userSession.WsConn == wsc {
					userSession.WsConn = nil
				}
				log.Printf("Websocket: User %s disconnect", uId)
				break
			}

			// parse client message
			var cMsg ClientMessage
			err = json.Unmarshal(data, &cMsg)
			if err != nil {
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			dbe := domain.CreateDbExecDeps(ctx, q, pool)

			// condition handlers
			if cMsg.Code == Ping {
				handlePingMessage(userSession, app)
			}
			if cMsg.Code == CaroJoinQueue {
				handleCaroJoinQueueMessage(userSession, app, dbe)
			}
			if cMsg.Code == CaroLeaveQueue {
				handleCaroLeaveQueueMessage(userSession, app)
			}
			if cMsg.Code == CaroPlayMove {
				var msgMove CaroPlayMoveMessageData
				err := json.Unmarshal(cMsg.Data, &msgMove)
				if err != nil {
					userSession.WsConn.WriteJSON(BuildServerMessage(Error, *userSession, app))
				} else {
					handleCaroPlayMoveMessage(userSession, app, msgMove, dbe)
				}

			}
			cancel()
		}
	}))
}

func BuildServerMessage(code ServerMessageCode, s domain.UserSession, a *domain.AppState) ServerMessage {
	sMsg := ServerMessage{Code: code, State: s.State}
	if sMsg.State == domain.PlayingCaro {
		match, existed := a.GetCaroMatch(s.CurrentMatchId)
		if existed == true {
			sMsg.Data = map[string]any{
				"currentMatch": match,
			}
		}

	}
	return sMsg
}

func handlePingMessage(s *domain.UserSession, a *domain.AppState) {
	s.WsConn.WriteJSON(BuildServerMessage(Ok, *s, a))
}

func handleCaroJoinQueueMessage(s *domain.UserSession, a *domain.AppState, dbe *domain.DbExecDeps) {
	match, err := a.UserJoinCaroQueue(dbe, s)
	if err != nil {
		s.WsConn.WriteJSON(BuildServerMessage(Error, *s, a))
		return
	}
	if match == nil {
		s.WsConn.WriteJSON(BuildServerMessage(Ok, *s, a))
		return
	}

	xSess, existed := a.GetUserSession(match.XPlayerId)
	if existed {
		xSess.WsConn.WriteJSON(BuildServerMessage(CaroMatchFound, *xSess, a))
	}

	oSess, existed := a.GetUserSession(match.OPlayerId)
	if existed {
		oSess.WsConn.WriteJSON(BuildServerMessage(CaroMatchFound, *oSess, a))
	}
}

func handleCaroLeaveQueueMessage(s *domain.UserSession, a *domain.AppState) {
	err := a.UserLeaveCaroQueue(s)
	if err != nil {
		s.WsConn.WriteJSON(BuildServerMessage(CaroMatchFound, *s, a))
		return
	}
	s.WsConn.WriteJSON(BuildServerMessage(CaroMatchFound, *s, a))
}

type CaroPlayMoveMessageData struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func handleCaroPlayMoveMessage(
	s *domain.UserSession,
	a *domain.AppState,
	move CaroPlayMoveMessageData,
	dbe *domain.DbExecDeps,
) {
	if s.State != domain.PlayingCaro {
		s.WsConn.WriteJSON(BuildServerMessage(Error, *s, a))
		return
	}

	match, existed := a.GetCaroMatch(s.CurrentMatchId)
	if !existed {
		s.WsConn.WriteJSON(BuildServerMessage(Error, *s, a))
		return
	}

	userPiece := domain.X
	if match.OPlayerId == s.UserId {
		userPiece = domain.O
	}

	err := match.Move(userPiece, move.X, move.Y)
	if err != nil {
		s.WsConn.WriteJSON(BuildServerMessage(Error, *s, a))
		return
	}

	xSess, xExisted := a.GetUserSession(match.XPlayerId)
	oSess, oExisted := a.GetUserSession(match.OPlayerId)

	if !xExisted || !oExisted {
		return
	}

	if !match.IsEnded {
		xSess.WsConn.WriteJSON(BuildServerMessage(CaroNewBoardState, *xSess, a))
		oSess.WsConn.WriteJSON(BuildServerMessage(CaroNewBoardState, *oSess, a))
		return
	}

	// match ended
	matchResult, err := a.ProcessCaroMatchEnded(dbe, *match)

	xMsg := BuildServerMessage(CaroMatchEnded, *xSess, a)
	xMsg.Data = map[string]any{
		"endedMatch": matchResult,
	}
	xSess.WsConn.WriteJSON(xMsg)

	oMsg := BuildServerMessage(CaroMatchEnded, *oSess, a)
	oMsg.Data = map[string]any{
		"endedMatch": matchResult,
	}
	oSess.WsConn.WriteJSON(oMsg)
}
