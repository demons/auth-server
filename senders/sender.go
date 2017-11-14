package senders

type Sender interface {
	Send(string, []byte) error
}
