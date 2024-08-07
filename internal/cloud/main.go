package cloud

type CloudProvider interface {
	GetStatus() (*WorkstationStatus, error)
	InitWorkstation(params *WorkstationInitParams) error
	StartWorkstation(params *WorkstationStartParams) error
	StopWorkstation(params *WorkstationStopParams) error
	DeleteWorkstation(params *WorkstationDeleteParams) error
}

type WorkstationInitParams struct {
	SshPubKey string
}

type WorkstationStartParams struct {
	SshPubKey string
}

type WorkstationSaveParams struct {
}

type WorkstationStopParams struct {
}

type WorkstationDeleteParams struct {
}

type WorkstationConnectParams struct {
}

type WorkstationStatus struct {
	IsActive  bool
	ID        int
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
	IPv4      string
}
