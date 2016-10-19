package entity

type CreateRequest struct {
	UserName    string    `json:"userName"`
	ClusterName string    `json:"clusterName"`
	Timeout     string    `json:"timeout"`
	SshUser     string    `json:"sshUser"`
	SkipInstall bool      `json:"skipInstall"`
	AddNodes    NodesInfo `json:"addNodes"`
}

type NodesInfo struct {
	PrivateKey     string   `json:"privateKey"`
	PrivateNicName string   `json:"privateNicName"`
	MasterNodes    []string `json:"masterNodes"`
	SalveNodes     []string `json:"salveNodes"`
}

type AddNodeRequest struct {
	UserName       string `json:"userName"`
	ClusterName    string `json:"clusterName"`
	SlaveType      string `json:"slaveType"`
	PrivateNicName string `json:"privateNicName"`
	Nodes          []Node `json:"node"`
}

type Node struct {
	Ip      string `json:"ip"`
	SshUser string `json:"sshUser"`
}

//type NodesInfo struct {
//	PrivateKey     string `json:"privateKey"`
//	PrivateNicName string `json:"privateNicName"`
//	MasterNodes    []Node `json:"masterNodes"`
//	SalveNodes     []Node `json:"salveNodes"`
//}
