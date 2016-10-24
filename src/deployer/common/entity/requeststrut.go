package entity

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
