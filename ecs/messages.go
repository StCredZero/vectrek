package ecs

import "github.com/StCredZero/vectrek/ecstypes"

type Sender interface {
	Send(msg ComponentMessage)
}

type Receiver interface {
	Receive() (ComponentMessage, bool)
}
type ComponentMessage struct {
	Entity  ecstypes.EntityID
	Payload any
}

type Pipe struct {
	Inbox chan ComponentMessage
}

func NewPipe() *Pipe {
	return &Pipe{
		Inbox: make(chan ComponentMessage, 1000),
	}
}

func (p *Pipe) Send(msg ComponentMessage) {
	p.Inbox <- msg
}

func (p *Pipe) Receive() (ComponentMessage, bool) {
	select {
	case msg := <-p.Inbox:
		return msg, true
	default:
		return ComponentMessage{}, false
	}
}
