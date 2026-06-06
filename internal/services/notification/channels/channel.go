package channels

type Channel string

var (
	Email     Channel = "email"
	Websocket Channel = "websocket"
	Push      Channel = "push"
	Database  Channel = "database"
)

type SendNotificationPayload struct {
	UserID   string
	Email    string
	Title    string
	Message  string
	Type     string
	Data     map[string]any
	Channels []Channel
}

type IChannel interface {
	Type() string
	Send(payload SendNotificationPayload) error
}
