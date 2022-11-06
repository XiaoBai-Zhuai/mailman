package entity

type ClientType int64

const (
	Producer ClientType = 1
	Consumer ClientType = 2

	NetWork = "tcp"
)

type MsgBody struct {
	Topic   string `json:"topic"`
	Message []byte `json:"message"`
}

type ClientMsg struct {
	Msg        MsgBody
	ClientType ClientType `json:"clientType"`
}
