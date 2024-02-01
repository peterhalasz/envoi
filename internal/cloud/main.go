package cloud

type CloudProvider interface {
	GetStatus() WorkstationStatus
}

type WorkstationStatus struct {
	Name string
}
