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
)

type ServerMessage struct {
	Code    ServerMessageCode  `json:"code"`
	Session domain.UserSession `json:"session"`
}

type ServerMessageCode string

const (
	Ok    ServerMessageCode = "OK"
	Error ServerMessageCode = "ERROR"
)

func InitWsRoute(a *fiber.App, q *db.Queries, conf config.Config, app *domain.AppState) {
	a.Get("/ws", wsAuthMiddleware(q, conf), websocket.New(func(wsc *websocket.Conn) {
		u := wsc.Locals("currentUser").(db.User)

		// save userSession
		userSession := app.SafeGetOrCreateUserSession(u.ID)
		if userSession.WsConn != nil {
			userSession.WsConn.Close()
		}
		userSession.WsConn = wsc
		log.Printf("Websocket: User %s connect", u.ID)

		for {
			_, data, err := wsc.ReadMessage()
			// disconnect
			if err != nil {
				wsc.Close()
				if userSession.WsConn == wsc {
					userSession.WsConn = nil
				}
				log.Printf("Websocket: User %s disconnect", u.ID)
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
				log.Println(app.CaroQueue)
			}
			if cMsg.Code == CaroLeaveQueue {
				handleCaroLeaveQueueMessage(userSession, app)
				log.Println(app.CaroQueue)
			}
		}
	}))
}

func handlePingMessage(s *domain.UserSession) {
	s.WsConn.WriteJSON(ServerMessage{Code: Ok, Session: *s})
}

func handleCaroJoinQueueMessage(s *domain.UserSession, a *domain.AppState) {
	err := domain.UserJoinCaroQueue(s, a)
	if err != nil {
		s.WsConn.WriteJSON(ServerMessage{Code: Error, Session: *s})
		return
	}
	s.WsConn.WriteJSON(ServerMessage{Code: Ok, Session: *s})
}

func handleCaroLeaveQueueMessage(s *domain.UserSession, a *domain.AppState) {
	err := domain.UserLeaveCaroQueue(s, a)
	if err != nil {
		s.WsConn.WriteJSON(ServerMessage{Code: Error, Session: *s})
		return
	}
	s.WsConn.WriteJSON(ServerMessage{Code: Ok, Session: *s})
}
