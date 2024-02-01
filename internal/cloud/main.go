package cloud

type CloudProvider interface {
	GetStatus() (*WorkstationStatus, error)
	InitWorkstation(params *WorkstationInitParams) error
	StartWorkstation(params *WorkstationStartParams) error
	SaveWorkstation(params *WorkstationSaveParams) error
	StopWorkstation(params *WorkstationStopParams) error
	DeleteWorkstation(params *WorkstationDeleteParams) error
}

type WorkstationInitParams struct {
}

type WorkstationStartParams struct {
}

type WorkstationSaveParams struct {
}

type WorkstationStopParams struct {
}

type WorkstationDeleteParams struct {
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
