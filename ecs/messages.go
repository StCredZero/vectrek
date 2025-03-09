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

type PipeConnectionManager struct {
	Pipe     *Pipe
	Receiver ecstypes.Receiver
	Sender   ecstypes.Sender
}

func NewPipeConnectionManager() *PipeConnectionManager {
	return &PipeConnectionManager{}
}

func (pcm *PipeConnectionManager) GetSender() ecstypes.Sender {
	return pcm.Sender
}

func (pcm *PipeConnectionManager) SetSender(sender ecstypes.Sender) {
	pcm.Sender = sender
}

func (pcm *PipeConnectionManager) GetReceiver() ecstypes.Receiver {
	return pcm.Receiver
}

func (pcm *PipeConnectionManager) SetReceiver(receiver ecstypes.Receiver) {
	pcm.Receiver = receiver
}
