package cloud

type CloudProvider interface {
	InitWorkstation(params WorkstationInitParams)
	GetStatus() (*WorkstationStatus, error)
}

type WorkstationInitParams struct {
}

type WorkstationStatus struct {
	IsActive  bool
	Name      string
	Memory    int
	Cpus      int
	Disk      int
	Region    string
	Image     string
	Size      string
	Status    string
	CreatedAt string
	Volume    string
}
