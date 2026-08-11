package ws

type ServerMessage struct {
	State UserState `json:"state"`
}