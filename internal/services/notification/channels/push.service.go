package channels

type PushChannel struct {
}

func NewPushChannel() *PushChannel {
	return &PushChannel{}
}

func (channel *PushChannel) Type() string {
	return string(Push)
}

func (channel *PushChannel) Send() error {
	return nil
}
