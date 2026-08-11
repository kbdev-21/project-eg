package ws

import (
	"backend/src/db"
	"encoding/json"
	"fmt"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
)

// TODO: dùng mutex
var connMap = map[pgtype.UUID]*websocket.Conn{}

func InitWsRoutes(a *fiber.App, q *db.Queries) {
	a.Get("/ws", wsMiddleware(q), websocket.New(func(wsc *websocket.Conn) {
		u := wsc.Locals("currentUser").(db.User)
		str := "You are " + u.Name
		wsc.WriteMessage(websocket.TextMessage, []byte(str))

		connMap[u.ID] = wsc

		for {	
			_, data, err := wsc.ReadMessage()
			// disconnect
			if err != nil {
				delete(connMap, u.ID)
				break
			}

			var cMsg ClientMessage
			err = json.Unmarshal(data, &cMsg)
			if err != nil {
				sendMessage(wsc, "Error")
				continue
			}

			if cMsg.Code == Ping {
				clientPingMessageHandler(u)
			}
		}
	}))
}

func sendMessage(wsc *websocket.Conn, msg string) {
	wsc.WriteMessage(websocket.TextMessage, []byte(msg))
}

func clientPingMessageHandler(u db.User) {
	fmt.Println("==========")
	for key, val := range connMap {
		fmt.Println(key)
		sendMessage(val, u.Name + " ping the server")
	}
}