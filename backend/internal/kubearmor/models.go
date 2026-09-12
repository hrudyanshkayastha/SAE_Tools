package kubearmor

type KubeArmorAlert struct {
	Timestamp         int64             `json:"Timestamp,omitempty"`
	UpdatedTime       string            `json:"UpdatedTime,omitempty"`
	ClusterName       string            `json:"ClusterName,omitempty"`
	HostName          string            `json:"HostName,omitempty"`
	NamespaceName     string            `json:"NamespaceName,omitempty"`
	PodName           string            `json:"PodName,omitempty"`
	Labels            string            `json:"Labels,omitempty"`
	ContainerID       string            `json:"ContainerID,omitempty"`
	ContainerName     string            `json:"ContainerName,omitempty"`
	ContainerImage    string            `json:"ContainerImage,omitempty"`
	HostPPID          int32             `json:"HostPPID,omitempty"`
	HostPID           int32             `json:"HostPID,omitempty"`
	PPID              int32             `json:"PPID,omitempty"`
	PID               int32             `json:"PID,omitempty"`
	UID               int32             `json:"UID,omitempty"`
	ParentProcessName string            `json:"ParentProcessName,omitempty"`
	ProcessName       string            `json:"ProcessName,omitempty"`
	PolicyName        string            `json:"PolicyName,omitempty"`
	Severity          string            `json:"Severity,omitempty"`
	Tags              string            `json:"Tags,omitempty"`
	Message           string            `json:"Message,omitempty"`
	Type              string            `json:"Type,omitempty"`
	Source            string            `json:"Source,omitempty"`
	Operation         string            `json:"Operation,omitempty"`
	Resource          string            `json:"Resource,omitempty"`
	Data              string            `json:"Data,omitempty"`
	EventData         map[string]string `json:"EventData,omitempty"`
	Enforcer          string            `json:"Enforcer,omitempty"`
	Action            string            `json:"Action,omitempty"`
	Result            string            `json:"Result,omitempty"`
}
