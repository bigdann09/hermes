package channels

type Channel string

var (
	Email     Channel = "email"
	Websocket Channel = "websocket"
	Push      Channel = "push"
	Database  Channel = "database"
)

type IChannel interface {
	Type() string
	Send() error
}
