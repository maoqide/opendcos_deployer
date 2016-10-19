package services

import (
	"deployer/common"
	"deployer/common/entity"
	"errors"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
)

var (
	BASE_PATH = "/opendcos/clusters/"

	DEPLOY_ERROR_INVALIDATE_CLUSTERNAME string = "INVALIDATE_CLUSTERNAME"
)

func cleanup(clusterDir string) {
	logrus.Infof("start cleanup...")
	//_, _, err = common.ExecCommand("sudo rm -rf " + clusterDir)
}

func CreateCluster(request entity.CreateRequest) (err error) {

	logrus.Infof("start createCluster...")
	logrus.Infof("createRequset: %v", request)

	clusterName := request.ClusterName
	username := request.UserName
	nodesInfo := request.AddNodes

	//check if username&clusterName validate
	if !checkName(username, clusterName) {
		logrus.Errorf("CreateCluster, invalidate clusterName.")
		err = errors.New(DEPLOY_ERROR_INVALIDATE_CLUSTERNAME)
		return
	}

	//generate clusterDir
	clusterDir := BASE_PATH + username + "/" + clusterName + "/"

	//preparation
	err = preparation(clusterName, clusterDir, request.Timeout, request.SshUser, nodesInfo.SalveNodes, nodesInfo.MasterNodes, nodesInfo.PrivateKey, nodesInfo.PrivateNicName)
	if err != nil {
		logrus.Errorf("createCluster, preparation failed. err is %v", err)
		return
	}

	//preCheck
	err = preCheck(clusterName, clusterDir, request.SkipInstall)
	if err != nil {
		logrus.Errorf("createCluster, preCheck failed. err is %v", err)
		return
	}

	//provision
	err = provision(clusterName, clusterDir)
	if err != nil {
		logrus.Errorf("createCluster, provision failed. err is %v", err)
		return
	}

	//postAction
	err = postAction(clusterName, clusterDir)
	if err != nil {
		logrus.Errorf("createCluster, postAction failed. err is %v", err)
		return
	}

	//backup
	err = backup()
	if err != nil {
		logrus.Errorf("createCluster, backup failed. err is %v", err)
		return
	}

	return
}

//prepare for deploy, create cluster directory, and execute --genconf
func preparation(clusterName string, clusterDir string, timeout string, sshUser string, slaves []string, masters []string, privateKey string, privateNicName string) (err error) {

	logrus.Infof("start preparation... clusterName: %s, clusterDir: %s, timeout: %s, sshUser: %s, slaves: %v, masters: %v, privateKey: %s, privateNicName: %s",
		clusterName, clusterDir, timeout, sshUser, slaves, masters, privateKey, privateNicName)

	_, _, err = common.ExecCommand("sudo mkdir -p " + clusterDir)
	if err != nil {
		logrus.Errorf("Preparation, ExecCommand err: %v", err)
		return
	}
	_, _, err = common.ExecCommand("sudo mkdir -p " + clusterDir + "genconf")
	if err != nil {
		logrus.Errorf("Preparation, ExecCommand err: %v", err)
		return
	}

	genConf(clusterDir+"genconf/", clusterName, timeout, sshUser, slaves, masters)
	genIPDetect(clusterDir+"genconf/", privateNicName)
	genSshKey(clusterDir+"genconf/", privateKey)

	//execute --genconf
	output, errput, err := common.ExecCommand("sudo bash script/exec.sh " + clusterDir + " genconf")
	if err != nil {
		logrus.Errorf("Preparation, ExecCommand err: %v", err)
		logrus.Infof("Preparation for cluster %s, errput: %s", clusterName, errput)
		//cleanup()
		return
	}
	logrus.Infof("Preparation for cluster %s, output: %s", clusterName, output)
	return
}

//check config to ensure the cluster can be deployed, execute --install-prereqs & --preflight
//not execute --install-prereqs if skipInstall is true
func preCheck(clusterName string, clusterDir string, skipInstall bool) (err error) {

	logrus.Infof("start preCheck... clusterName: %s, clusterDir: %s, skipInstall: %t", clusterName, clusterDir, skipInstall)

	//skipInstall is false, execute --install-prereqs
	if !skipInstall {
		output, errput, err1 := common.ExecCommand("sudo bash script/exec.sh " + clusterDir + " install-prereqs")
		if err1 != nil {
			logrus.Errorf("preCheck --install-prereqs, ExecCommand err: %v", err)
			logrus.Infof("preCheck --install-prereqs for cluster %s, errput: %s", clusterName, errput)
			//cleanup()
			return err1
		}
		logrus.Infof("preCheck --install-prereqs for cluster %s, output: %s", clusterName, output)
	}

	//execute --preflight
	output, errput, err := common.ExecCommand("sudo bash script/exec.sh " + clusterDir + " preflight")
	if err != nil {
		logrus.Errorf("preCheck --preflight, ExecCommand err: %v", err)
		logrus.Infof("preCheck --preflight for cluster %s, errput: %s", clusterName, errput)
		//cleanup()
		return
	}
	logrus.Infof("preCheck --preflight for cluster %s, output: %s", clusterName, output)

	return
}

//deploy cluster, execute --deploy
func provision(clusterName string, clusterDir string) (err error) {

	logrus.Infof("start provision... clusterName: %s, clusterDir: %s", clusterName, clusterDir)

	//execute --deploy
	output, errput, err := common.ExecCommand("sudo bash script/exec.sh " + clusterDir + " deploy")
	if err != nil {
		logrus.Errorf("provision --deploy, ExecCommand err: %v", err)
		logrus.Infof("provision --deploy for cluster %s, errput: %s", clusterName, errput)
		//cleanup()
		return
	}
	logrus.Infof("provision --deploy for cluster %s, output: %s", clusterName, output)
	return
}

//postaction of deployment, exec --postflight
func postAction(clusterName string, clusterDir string) (err error) {

	logrus.Infof("start postAction... clusterName: %s, clusterDir: %s", clusterName, clusterDir)

	//execute --postflight
	output, errput, err := common.ExecCommand("sudo bash script/exec.sh " + clusterDir + " postflight")
	if err != nil {
		logrus.Errorf("postAction --postflight, ExecCommand err: %v", err)
		logrus.Infof("postAction --postflight for cluster %s, errput: %s", clusterName, errput)
		//cleanup()
		return
	}
	logrus.Infof("postAction --postflight for cluster %s, output: %s", clusterName, output)
	return
}

//backup dcos-install.tar
func backup() (err error) {
	return
}

//check if username & clusterName validate
func checkName(username string, clusterName string) (validate bool) {

	//check if existed

	validate = true
	return
}

//generate config.yaml
//path should be absolute and end with "/"
func genConf(path string, clusterName string, timeout string, sshUser string, slaves []string, masters []string) (err error) {

	logrus.Infof("start create file config.yaml... path: %s, clusterName: %s, timeout: %s, sshUser: %s, slaves: %v, masters: %v",
		path, clusterName, timeout, sshUser, slaves, masters)

	var strSlaves = ""
	var strMasters = ""
	for _, sip := range slaves {
		strSlaves = strSlaves + `- ` + sip + "\n"
	}
	for _, mip := range masters {
		strMasters = strMasters + `- ` + mip + "\n"
	}
	logrus.Infof("genConf, slaveStr: %s", strSlaves)
	logrus.Infof("genConf, masterStr: %s", strMasters)

	fileStr := `---
agent_list:
` + strSlaves + `bootstrap_url: file:///opt/dcos_install_tmp
cluster_name: ` + clusterName + `
exhibitor_storage_backend: static
ip_detect_filename: /genconf/ip-detect
master_discovery: static
master_list:
` + strMasters + `process_timeout: ` + timeout + `
resolvers:
- 8.8.8.8
- 114.114.114.114
ssh_port: 22
ssh_user: ` + sshUser + `
`

	logrus.Debugf("genConf, config.yaml: ", fileStr)
	err = ioutil.WriteFile(path+"config.yaml", []byte(fileStr), 0644)
	if err != nil {
		logrus.Errorf("genConf failed, err is %v", err)
	}
	logrus.Infof("file config.yaml created.")
	return
}

//generate ip-detect
//path should be absolute and end with "/"
func genIPDetect(path string, privateNicName string) (err error) {

	logrus.Infof("start create file ip-detect... path: %s", path)

	fileStr := `#!/usr/bin/env bash
set -o nounset -o errexit
export PATH=/usr/sbin:/usr/bin:$PATH
echo $(ip addr show ` + privateNicName + ` | grep -Eo '[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}' | head -1)
`

	logrus.Debugf("genIPDetect, ip-detect: ", fileStr)
	err = ioutil.WriteFile(path+"ip-detect", []byte(fileStr), 0644)
	if err != nil {
		logrus.Errorf("genIPDetect failed, err is %v", err)
	}
	logrus.Infof("file ip-detect created.")
	return
}

//generate ssh_key
//path should be absolute and end with "/"
func genSshKey(path string, privateKey string) (err error) {

	logrus.Infof("start create file ssh_key... path: %s, privateKey: %s", path, privateKey)

	decodeKey, err := common.Base64Decode([]byte(privateKey))
	if err != nil {
		logrus.Errorf("decode private key error %v", err)
		return
	}

	err = ioutil.WriteFile(path+"ssh_key", []byte(privateKey), 0600)
	if err != nil {
		logrus.Errorf("genSshKey failed, err is %v", err)
	}
	logrus.Infof("file ssh_key created.")
	return
}
