package router

import (
	"backend/src/config"
	"backend/src/db"
	"backend/src/domain"
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
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
	Ok                ServerMessageCode = "OK"
	Error             ServerMessageCode = "ERROR"
	CaroMatchFound    ServerMessageCode = "CARO:MATCH_FOUND"
	CaroNewBoardState ServerMessageCode = "CARO:NEW_BOARD_STATE"
	CaroMatchEnded    ServerMessageCode = "CARO:MATCH_ENDED"
)

func InitWsRoute(a *fiber.App, q *db.Queries, conf config.Config, app *domain.AppState) {
	a.Get("/ws", wsAuthMiddleware(q, conf), websocket.New(func(wsc *websocket.Conn) {
		uId := wsc.Locals("currentUser").(db.User).ID

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

			// condition handlers
			if cMsg.Code == Ping {
				handlePingMessage(userSession, app)
			}
			if cMsg.Code == CaroJoinQueue {
				handleCaroJoinQueueMessage(userSession, app)
			}
			if cMsg.Code == CaroLeaveQueue {
				handleCaroLeaveQueueMessage(userSession, app)
			}
			if cMsg.Code == CaroPlayMove {
				handleCaroPlayMoveMessage(userSession, app, cMsg.Data)
			}
		}
	}))
}

func buildServerMessage(code ServerMessageCode, s domain.UserSession, a *domain.AppState) ServerMessage {
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
	s.WsConn.WriteJSON(buildServerMessage(Ok, *s, a))
}

func handleCaroJoinQueueMessage(s *domain.UserSession, a *domain.AppState) {
	match, err := a.UserJoinCaroQueue(s)
	if err != nil {
		s.WsConn.WriteJSON(buildServerMessage(Error, *s, a))
		return
	}
	if match == nil {
		s.WsConn.WriteJSON(buildServerMessage(Ok, *s, a))
		return
	}

	xSess, existed := a.GetUserSession(match.XPlayerId)
	if existed {
		xSess.WsConn.WriteJSON(buildServerMessage(CaroMatchFound, *xSess, a))
	}

	oSess, existed := a.GetUserSession(match.OPlayerId)
	if existed {
		oSess.WsConn.WriteJSON(buildServerMessage(CaroMatchFound, *oSess, a))
	}
}

func handleCaroLeaveQueueMessage(s *domain.UserSession, a *domain.AppState) {
	err := a.UserLeaveCaroQueue(s)
	if err != nil {
		s.WsConn.WriteJSON(buildServerMessage(CaroMatchFound, *s, a))
		return
	}
	s.WsConn.WriteJSON(buildServerMessage(CaroMatchFound, *s, a))
}

type CaroPlayMoveMessageData struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func handleCaroPlayMoveMessage(s *domain.UserSession, a *domain.AppState, msgData json.RawMessage) {
	if s.State != domain.PlayingCaro {
		s.WsConn.WriteJSON(buildServerMessage(Error, *s, a))
		return
	}

	var msgMove CaroPlayMoveMessageData
	err := json.Unmarshal(msgData, &msgMove)
	if err != nil {
		s.WsConn.WriteJSON(buildServerMessage(Error, *s, a))
		return
	}

	match, existed := a.GetCaroMatch(s.CurrentMatchId)
	if !existed {
		s.WsConn.WriteJSON(buildServerMessage(Error, *s, a))
		return
	}

	userPiece := domain.X
	if match.OPlayerId == s.UserId {
		userPiece = domain.O
	}

	err = match.Move(userPiece, msgMove.X, msgMove.Y)
	if err != nil {
		s.WsConn.WriteJSON(buildServerMessage(Error, *s, a))
		return
	}

	xSess, xExisted := a.GetUserSession(match.XPlayerId)
	oSess, oExisted := a.GetUserSession(match.OPlayerId)

	if !xExisted || !oExisted {
		return
	}

	if !match.IsEnded {
		xSess.WsConn.WriteJSON(buildServerMessage(CaroNewBoardState, *xSess, a))
		oSess.WsConn.WriteJSON(buildServerMessage(CaroNewBoardState, *oSess, a))
		return
	}

	// TODO: impl post match correctly later
	matchResult := *match
	xSess.State = domain.Idle
	xSess.CurrentMatchId = uuid.Nil
	oSess.State = domain.Idle
	oSess.CurrentMatchId = uuid.Nil

	xMsg := buildServerMessage(CaroMatchEnded, *xSess, a)
	xMsg.Data = map[string]any{
		"endedMatch": matchResult,
	}
	xSess.WsConn.WriteJSON(xMsg)

	oMsg := buildServerMessage(CaroMatchEnded, *oSess, a)
	oMsg.Data = map[string]any{
		"endedMatch": matchResult,
	}
	oSess.WsConn.WriteJSON(oMsg)
}
