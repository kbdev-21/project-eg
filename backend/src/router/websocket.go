package router

import (
	"backend/src/config"
	"backend/src/db"
	"backend/src/domain"
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

type ClientMessage struct {
	Code ClientMessageCode `json:"code"`
	Data any               `json:"data"`
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
	Code    ServerMessageCode  `json:"code"`
	Session domain.UserSession `json:"session"`
	Data    any                `json:"data"`
}

type ServerMessageCode string

const (
	Ok         ServerMessageCode = "OK"
	Error      ServerMessageCode = "ERROR"
	MatchFound ServerMessageCode = "MATCH_FOUND"
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
				handlePingMessage(userSession)
			}
			if cMsg.Code == CaroJoinQueue {
				handleCaroJoinQueueMessage(userSession, app)
			}
			if cMsg.Code == CaroLeaveQueue {
				handleCaroLeaveQueueMessage(userSession, app)
			}
		}
	}))
}

func handlePingMessage(s *domain.UserSession) {
	s.WsConn.WriteJSON(ServerMessage{Code: Ok, Session: *s})
}

func handleCaroJoinQueueMessage(s *domain.UserSession, a *domain.AppState) {
	match, err := a.UserJoinCaroQueue(s)
	if err != nil {
		s.WsConn.WriteJSON(ServerMessage{Code: Error, Session: *s})
		return
	}
	if match == nil {
		s.WsConn.WriteJSON(ServerMessage{Code: Ok, Session: *s})
		return
	}

	xSess, existed := a.GetUserSession(match.XPlayerId)
	if existed {
		xSess.WsConn.WriteJSON(ServerMessage{
			Code:    MatchFound,
			Session: *xSess,
			Data:    *match,
		})
	}

	oSess, existed := a.GetUserSession(match.OPlayerId)
	if existed {
		oSess.WsConn.WriteJSON(ServerMessage{
			Code:    MatchFound,
			Session: *oSess,
			Data:    *match,
		})
	}
}

func handleCaroLeaveQueueMessage(s *domain.UserSession, a *domain.AppState) {
	err := a.UserLeaveCaroQueue(s)
	if err != nil {
		s.WsConn.WriteJSON(ServerMessage{Code: Error, Session: *s})
		return
	}
	s.WsConn.WriteJSON(ServerMessage{Code: Ok, Session: *s})
}
