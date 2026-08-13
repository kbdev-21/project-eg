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
	Ping      ClientMessageCode = "PING"
	JoinQueue ClientMessageCode = "JOIN_QUEUE"
)

func InitWsRoute(a *fiber.App, q *db.Queries, conf config.Config, aState *domain.AppState) {
	a.Get("/ws", wsAuthMiddleware(q, conf), websocket.New(func(wsc *websocket.Conn) {
		u := wsc.Locals("currentUser").(db.User)

		session := aState.GetOrCreateUserSession(u.ID)
		if session.WsConn != nil {
			session.WsConn.Close()
		}
		session.WsConn = wsc

		log.Printf("Websocket: User %s connect", u.ID)

		for {
			_, data, err := wsc.ReadMessage()
			// disconnect
			if err != nil {
				wsc.Close()
				if session.WsConn == wsc {
					session.WsConn = nil
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
				handlePingMessage(session)
			}
		}
	}))
}

func handlePingMessage(session *domain.UserSession) {
	encoded, err := json.Marshal(session)
	if err != nil {
		session.WsConn.WriteMessage(websocket.TextMessage, []byte("Some thing going wrong"))
		return
	}
	session.WsConn.WriteMessage(websocket.TextMessage, encoded)
}
