package entity

//type CreateRequest struct {
//	UserName    string    `json:"userName"`
//	ClusterName string    `json:"clusterName"`
//	Timeout     string    `json:"timeout"`
//	SshUser     string    `json:"sshUser"`
//	SkipInstall bool      `json:"skipInstall"`
//	AddNodes    NodesInfo `json:"addNodes"`
//}

//type NodesInfo struct {
//	PrivateKey     string   `json:"privateKey"`
//	PrivateNicName string   `json:"privateNicName"`
//	MasterNodes    []string `json:"masterNodes"`
//	SalveNodes     []string `json:"salveNodes"`
//}

type CreateRequest struct {
	UserName       string     `json:"userName"`
	SkipInstall    bool       `json:"skipInstall"`
	PrivateKey     string     `json:"privateKey"`
	PrivateNicName string     `json:"privateNicName"`
	Config         DCOSConfig `json:"config"`
}

type AddNodeRequest struct {
	UserName       string   `json:"userName"`
	ClusterName    string   `json:"clusterName"`
	SlaveType      string   `json:"slaveType"` //'slave' or 'slave_public'
	PrivateNicName string   `json:"privateNicName"`
	SshUser        string   `json:"sshUser"`
	Nodes          []string `json:"nodes"`
}
