package ws

import (
	"backend/src/db"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

func InitWsRoutes(a *fiber.App, q *db.Queries) {
	a.Get("/ws", wsMiddleware, websocket.New(func(wsc *websocket.Conn) {
		for {
			_, message, err := wsc.ReadMessage()
			if err != nil {
				break
			}

			err = wsc.WriteMessage(1, message)
			if err != nil {
				break
			}
		}
	}))
}
