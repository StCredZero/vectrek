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

type ComponentStorage interface {
	SystemID() SystemID
}

type Component interface {
	Init(sm SystemManager, entity EntityID) error
	Update(sm SystemManager) error
	SystemID() SystemID
}

type Sender interface {
	Send(msg ComponentMessage)
}

type Receiver interface {
	Receive() (ComponentMessage, bool)
}

type ComponentMessage struct {
	Entity  EntityID
	Payload any
}
