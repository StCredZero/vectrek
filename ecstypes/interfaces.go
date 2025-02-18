package ecstypes

type System interface {
	IsSystem()
	SystemID() SystemID
}

type SystemManager interface {
	GetSystem(id SystemID) (System, error)
	AddComponent(e EntityID, component Component) error
	GetComponent(systemID SystemID, e EntityID) (Component, bool)
}

type ComponentStorage interface {
	SystemID() SystemID
}

type Component interface {
	Init(sm SystemManager, entity EntityID) error
	Update() error
	SystemID() SystemID
}
