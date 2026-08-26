package router

import (
	"backend/src/config"
	"backend/src/domain"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
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

func InitWsRoute(fib *fiber.App, a *domain.AppState, cfg config.Config, hub *WsConnHub) {
	fib.Get("/ws", wsAuthMiddleware(a, cfg), websocket.New(func(wsc *websocket.Conn) {
		uId := wsc.Locals("currentUser").(domain.User).Id

		// save userSession
		userSession := a.GetOrCreateUserSession(uId)
		hub.Register(uId, wsc)
		log.Printf("WS: User %s connect", uId)

		for {
			_, data, err := wsc.ReadMessage()
			// disconnect
			if err != nil {
				wsc.Close()
				hub.Unregister(uId, wsc)
				log.Printf("WS: User %s disconnect", uId)
				break
			}

			// parse client message
			var cMsg ClientMessage
			err = json.Unmarshal(data, &cMsg)
			if err != nil {
				continue
			}

			// condition handlers
			log.Printf("WS: User %s send %s message", uId, cMsg.Code)
			if cMsg.Code == Ping {
				handlePingMessage(userSession, a, hub)
			}
			if cMsg.Code == CaroJoinQueue {
				handleCaroJoinQueueMessage(userSession, a, hub)
			}
			if cMsg.Code == CaroLeaveQueue {
				handleCaroLeaveQueueMessage(userSession, a, hub)
			}
			if cMsg.Code == CaroPlayMove {
				var msgMove CaroPlayMoveMessageData
				err := json.Unmarshal(cMsg.Data, &msgMove)
				if err != nil {
					hub.Send(userSession.UserId, BuildServerMessage(Error, *userSession, a))
					continue
				}
				handleCaroPlayMoveMessage(userSession, a, hub, msgMove)

			}
		}
	}))
}

func BuildServerMessage(code ServerMessageCode, s domain.UserSession, a *domain.AppState) ServerMessage {
	sMsg := ServerMessage{Code: code, State: s.State}
	if sMsg.State == domain.PlayingCaro {
		match, existed := a.GetCaroMatch(s.CurrentMatchId)
		if existed == true {
			sMsg.Data = map[string]any{
				"match": match,
			}
		}

	}
	return sMsg
}

func handlePingMessage(s *domain.UserSession, a *domain.AppState, hub *WsConnHub) {
	hub.Send(s.UserId, BuildServerMessage(Ok, *s, a))
}

func handleCaroJoinQueueMessage(s *domain.UserSession, a *domain.AppState, hub *WsConnHub) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	match, err := a.UserJoinCaroQueue(ctx, s)
	if err != nil {
		hub.Send(s.UserId, BuildServerMessage(Error, *s, a))
		return
	}
	if match == nil {
		hub.Send(s.UserId, BuildServerMessage(Ok, *s, a))
		return
	}

	xSess, existed := a.GetUserSession(match.XPlayerId)
	if existed {
		hub.Send(xSess.UserId, BuildServerMessage(CaroMatchFound, *xSess, a))
	}

	oSess, existed := a.GetUserSession(match.OPlayerId)
	if existed {
		hub.Send(oSess.UserId, BuildServerMessage(CaroMatchFound, *oSess, a))
	}
}

func handleCaroLeaveQueueMessage(s *domain.UserSession, a *domain.AppState, hub *WsConnHub) {
	err := a.UserLeaveCaroQueue(s)
	if err != nil {
		hub.Send(s.UserId, BuildServerMessage(CaroMatchFound, *s, a))
		return
	}
	hub.Send(s.UserId, BuildServerMessage(CaroMatchFound, *s, a))
}

type CaroPlayMoveMessageData struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func handleCaroPlayMoveMessage(
	s *domain.UserSession,
	a *domain.AppState,
	hub *WsConnHub,
	move CaroPlayMoveMessageData,
) {
	if s.State != domain.PlayingCaro {
		hub.Send(s.UserId, BuildServerMessage(Error, *s, a))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	match, existed := a.GetCaroMatch(s.CurrentMatchId)
	if !existed {
		hub.Send(s.UserId, BuildServerMessage(Error, *s, a))
		return
	}

	userPiece := domain.X
	if match.OPlayerId == s.UserId {
		userPiece = domain.O
	}

	err := match.Move(userPiece, move.X, move.Y)
	if err != nil {
		hub.Send(s.UserId, BuildServerMessage(Error, *s, a))
		return
	}

	xSess, xExisted := a.GetUserSession(match.XPlayerId)
	oSess, oExisted := a.GetUserSession(match.OPlayerId)

	if !xExisted || !oExisted {
		return
	}

	if match.Status == domain.Playing {
		hub.Send(xSess.UserId, BuildServerMessage(CaroNewBoardState, *xSess, a))
		hub.Send(oSess.UserId, BuildServerMessage(CaroNewBoardState, *oSess, a))
		return
	}

	// match ended
	endedMatch, err := a.ProcessCaroMatchEnded(ctx, *match)
	if err != nil {
		return
	}

	xMsg := BuildServerMessage(CaroMatchEnded, *xSess, a)
	xMsg.Data = map[string]any{
		"match": endedMatch,
	}
	hub.Send(xSess.UserId, xMsg)

	oMsg := BuildServerMessage(CaroMatchEnded, *oSess, a)
	oMsg.Data = map[string]any{
		"match": endedMatch,
	}
	hub.Send(oSess.UserId, oMsg)
}
