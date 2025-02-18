package ecs

import "github.com/StCredZero/vectrek/ecstypes"

type Pipe struct {
	Inbox chan ecstypes.ComponentMessage
}

func NewPipe() *Pipe {
	return &Pipe{
		Inbox: make(chan ecstypes.ComponentMessage, 1000),
	}
}

func (p *Pipe) Send(msg ecstypes.ComponentMessage) {
	p.Inbox <- msg
}

func (p *Pipe) Receive() (ecstypes.ComponentMessage, bool) {
	select {
	case msg := <-p.Inbox:
		return msg, true
	default:
		return ecstypes.ComponentMessage{}, false
	}
}
