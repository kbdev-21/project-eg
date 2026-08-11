package ws

type ClientMessageCode string

const (
	Ping      ClientMessageCode = "PING"
	JoinQueue ClientMessageCode = "JOIN_QUEUE"
)

type ClientMessage struct {
	Code ClientMessageCode `json:"code"`
	Data any               `json:"data"`
}

type PingMessageData struct {
}