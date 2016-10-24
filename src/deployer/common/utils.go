package common

import (
	"encoding/base64"
	"os"
)

const (
	base64Table = "ABCDEFGHIJKLMNOPQRSTpqrstuvwxyz0123456789+/UVWXYZabcdefghijklmno"
)

var coder = base64.NewEncoding(base64Table)

func Base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func Base64Decode(src []byte) ([]byte, error) {
	return coder.DecodeString(string(src))
}

//check if file or directory exist
func PathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getSshUserAndPort(nodeDir string)(sshUser string, sshPort string,err error){

	logrus.Infof("make file to storage ip, path is: %s", clusterDir)

	configString, err := ioutil.ReadFile(nodeDir)
	if err!=nil {
		logrus.Errorf("read file failed,filename is:%s,err is: %s", nodeDir, err)
	}
	
	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte((string)configString), &m)

    if err != nil {
          log.Fatalf("error: %v", err)
    }

	sshUser= yaml[ssh_user] 
	sshPort= yaml[ssh_port]

	return
}

//storage node ip
func StorageNodeIp(nodeIp string, clusterDir string)(err error){

	exist,err := PathExist(clusterDir+"nodeip")

	if exist{
		file := os.Open(clusterDir+"nodeip")
		defer file.Close()
		file.WriteString(","+nodeIp)
	}else{
		err = ioutil.WriteFile(clusterDir+"nodeip", []byte(nodeIp), 0644)
		if err != nil {

		}
	}
	return
}

//get node ip
func GetNodeIp(clusterDir string)( nodeIp []string, err error){
	nodeString, err := ioutil.ReadFile(clusterDir+"nodeip")
	nodeIp = strings.Split((string)nodeString, ",")
	return
}

//query ip whether exist in nodefile
func QueryNodeIp(clusterDir string, nodeIp string)(exist bool, err error){
	nodeIp[], err := GetNodeIp(clusterDir)
	for _, node := range nodeIp {
		if strings.EqualFold(node, nodeIp){
			exist = true
		}
	} 
	return
}
