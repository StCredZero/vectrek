package ecstypes

type System interface {
	IsSystem()
	SystemID() SystemID
	Iterate() []error
}

type SystemManager interface {
	GetSystem(id SystemID) (System, error)
	AddComponent(e EntityID, component Component) error
	GetComponent(systemID SystemID, e EntityID) (Component, bool)
	GetSender() Sender
	GetReceiver() Receiver
	GetCounter() uint64
	GetName() string
}

type Component interface {
	Init(sm SystemManager, entity EntityID) error
	SystemID() SystemID
}

type Sender interface {
	Send(msg ComponentMessage)
}

type Receiver interface {
	Receive() (ComponentMessage, bool)
}

type ConnectionManager interface {
	GetSender() Sender
	SetSender(sender Sender)
	GetReceiver() Receiver
	SetReceiver(receiver Receiver)
}

type ComponentMessage struct {
	Entity  EntityID
	Payload any
}
