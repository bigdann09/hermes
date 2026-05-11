package channels

type MailChannel struct {
}

func NewMailChannel() *MailChannel {
	return &MailChannel{}
}

func (channel *MailChannel) Type() string {
	return string(Email)
}

func (channel *MailChannel) Send() error {
	return nil
}
