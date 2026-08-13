package router

import (
	"backend/src/config"
	"backend/src/db"
	"backend/src/domain"
	"encoding/json"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
)

// TODO: dùng mutex
var sessionMap = map[pgtype.UUID]*websocket.Conn{}

type ClientMessage struct {
	Code ClientMessageCode `json:"code"`
	Data any               `json:"data"`
}

type ClientMessageCode string

const (
	Ping      ClientMessageCode = "PING"
	JoinQueue ClientMessageCode = "JOIN_QUEUE"
)

func InitWsRoute(a *fiber.App, q *db.Queries) {
	a.Get("/ws", wsAuthMiddleware(q), websocket.New(func(wsc *websocket.Conn) {
		u := wsc.Locals("currentUser").(db.User)
		str := "You are " + u.Name
		wsc.WriteMessage(websocket.TextMessage, []byte(str))

		sessionMap[u.ID] = wsc

		for {
			_, data, err := wsc.ReadMessage()
			// disconnect
			if err != nil {
				delete(sessionMap, u.ID)
				wsc.Close()
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
				// do Ping message handler
				wsc.WriteMessage(websocket.TextMessage, []byte("U ping the server!"))
			}
		}
	}))
}

func wsAuthMiddleware(q *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			token := c.Query("token", "")

			payload, err := domain.VerifyToken(token, config.AppJwkSet)
			if err != nil {
				return c.SendStatus(401)
			}

			currentUser, err := domain.SyncUserFromTokenPayload(c, q, payload)
			if err != nil {
				return c.SendStatus(401)
			}

			c.Locals("currentUser", currentUser)

			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}