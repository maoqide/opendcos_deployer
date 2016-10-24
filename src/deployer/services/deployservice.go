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

//remove the cluster directory
func cleanup(clusterDir string) (err error) {

	logrus.Infof("start cleanup... clusterDir: %s", clusterDir)

	_, errput, err := common.ExecCommand("sudo rm -rf " + clusterDir)
	if err != nil {
		logrus.Errorf("cleanup error, err is %s, errput: %s", err, errput)
	}

	logrus.Infof("cleanup clusterDir: %s, finished", clusterDir)
	return
}

//TODO
func cleanNode() (err error) {
	return
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
	clusterDir := genClusterDir(username, clusterName)

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
	err = backup(clusterDir)
	if err != nil {
		logrus.Errorf("createCluster, backup failed. err is %v", err)
		return
	}

	//remove tar package
	commandStr := "sudo rm -f " + clusterDir + "*.tar"
	logrus.Infof("CreateCluster, execute command: %s", commandStr)
	output, errput, err := common.ExecCommand(commandStr)
	if err != nil {
		logrus.Errorf("CreateCluster, ExecCommand err: %v", err)
		logrus.Infof("CreateCluster %s, errput: %s", clusterName, errput)
		return
	}
	logrus.Infof("CreateCluster %s, output: %s", clusterName, output)

	return
}

func DeleteCluster(username string, clusterName string) (err error) {

	logrus.Infof("start DeleteCluster... username: %s, clusterName: %s", username, clusterName)

	//generate clusterDir
	clusterDir := genClusterDir(username, clusterName)
	privateKeyPath := clusterDir + "genconf/ssh_key"

	//check if cluster exists
	exist, _ := common.PathExist(clusterDir)
	if !exist {
		logrus.Errorf("DeleteCluster, cluster not existed. clusterDir: %s", clusterDir)
		return
	}

	//execute --uninstall
	commandStr := "sudo bash script/exec.sh " + clusterDir + " uninstall"
	logrus.Infof("DeleteCluster,execute command: %s", commandStr)
	output, errput, err := common.ExecCommand(commandStr)
	if err != nil {
		logrus.Errorf("DeleteCluster, ExecCommand err: %v", err)
		logrus.Infof("DeleteCluster %s, errput: %s", clusterName, errput)
		return
	}
	logrus.Infof("DeleteCluster %s, output: %s", clusterName, output)

	//TODO
	//find slaves from cluster dir
	nodes := []string{}
	sshUser := "root"
	for _, nodeip := range nodes {
		go deleteSingleNode(nodeip, sshUser, privateKeyPath)
	}

	cleanup(clusterDir)

	return

}

func AddNodes(request entity.AddNodeRequest) (err error) {

	logrus.Infof("start AddNodes...")
	logrus.Infof("AddNodeRequest: %v", request)

	clusterName := request.ClusterName
	username := request.UserName
	nodes := request.Nodes

	//generate clusterDir
	clusterDir := genClusterDir(username, clusterName)
	privateKeyPath := clusterDir + "genconf/ssh_key"

	//check if cluster exists
	exist, _ := common.PathExist(clusterDir)
	if !exist {
		logrus.Errorf("AddNodes, cluster not existed. clusterDir: %s", clusterDir)
		return
	}

	//check if dcos-installer.tar exists
	exist, _ = common.PathExist(clusterDir + "genconf/serve/dcos-install.tar")
	if !exist {
		logrus.Infof("AddNodes, backup file not exists. try backup...")
		errb := backup(clusterDir)
		if errb != nil {
			logrus.Errorf("AddNodes, backup failed, err is %v", err)
			return
		}
	}

	//start to add nodes
	for _, nodeip := range nodes {

		go addSingleNode(nodeip, request.SshUser, privateKeyPath, clusterDir, request.SlaveType)

	}

	return
}

func DeleteNode(username string, clusterName string, ip string) (err error) {

	logrus.Infof("start DeleteNode...")
	sshUser := "root"

	//TODO
	//check

	//generate clusterDir
	clusterDir := genClusterDir(username, clusterName)
	privateKeyPath := clusterDir + "genconf/ssh_key"

	go deleteSingleNode(ip, sshUser, privateKeyPath)

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
	commandStr := "sudo bash script/exec.sh " + clusterDir + " genconf"
	logrus.Infof("preparation, execute command: %s", commandStr)
	output, errput, err := common.ExecCommand(commandStr)
	if err != nil {
		logrus.Errorf("Preparation, ExecCommand err: %v", err)
		logrus.Infof("Preparation for cluster %s, errput: %s", clusterName, errput)
		//cleanup()
		return
	}
	logrus.Infof("Preparation for cluster %s, output: %s", clusterName, output)
	logrus.Infof("Preparation finished clusterName: %s, clusterDir: %s", clusterName, clusterDir)

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
	commandStr := "sudo bash script/exec.sh " + clusterDir + " preflight"
	logrus.Infof("preCheck, execute command: %s", commandStr)
	output, errput, err := common.ExecCommand(commandStr)
	if err != nil {
		logrus.Errorf("preCheck --preflight, ExecCommand err: %v", err)
		logrus.Infof("preCheck --preflight for cluster %s, errput: %s", clusterName, errput)
		//cleanup()
		return
	}
	logrus.Infof("preCheck --preflight for cluster %s, output: %s", clusterName, output)
	logrus.Infof("preCheck finished. clusterName: %s, clusterDir: %s", clusterName, clusterDir)

	return
}

//deploy cluster, execute --deploy
func provision(clusterName string, clusterDir string) (err error) {

	logrus.Infof("start provision... clusterName: %s, clusterDir: %s", clusterName, clusterDir)

	//execute --deploy
	commandStr := "sudo bash script/exec.sh " + clusterDir + " deploy"
	logrus.Infof("provision, execute command: %s", commandStr)
	output, errput, err := common.ExecCommand(commandStr)
	if err != nil {
		logrus.Errorf("provision --deploy, ExecCommand err: %v", err)
		logrus.Infof("provision --deploy for cluster %s, errput: %s", clusterName, errput)
		//cleanup()
		return
	}
	logrus.Infof("provision --deploy for cluster %s, output: %s", clusterName, output)
	logrus.Infof(" provision finished. clusterName: %s, clusterDir: %s", clusterName, clusterDir)
	return
}

//postaction of deployment, exec --postflight
func postAction(clusterName string, clusterDir string) (err error) {

	logrus.Infof("start postAction... clusterName: %s, clusterDir: %s", clusterName, clusterDir)

	//execute --postflight
	commandStr := "sudo bash script/exec.sh " + clusterDir + " postflight"
	logrus.Infof("postAction, execute command: %s", commandStr)
	output, errput, err := common.ExecCommand(commandStr)
	if err != nil {
		logrus.Errorf("postAction --postflight, ExecCommand err: %v", err)
		logrus.Infof("postAction --postflight for cluster %s, errput: %s", clusterName, errput)
		//cleanup()
		return
	}
	logrus.Infof("postAction --postflight for cluster %s, output: %s", clusterName, output)
	logrus.Infof("postAction finished. clusterName: %s, clusterDir: %s", clusterName, clusterDir)
	return
}

//backup dcos-install.tar
func backup(clusterDir string) (err error) {

	logrus.Infof("start backup...  clusterDir: %s", clusterDir)

	//command: sudo tar cf $clusterDir/genconf/serve/dcos-install.tar -C $clusterDir/genconf/serve .
	commandStr := "sudo tar cf " + clusterDir + "genconf/serve/dcos-install.tar -C " + clusterDir + "genconf/serve ."
	logrus.Infof("backup, execute command: %s", commandStr)
	_, errput, err := common.ExecCommand(commandStr)
	if err != nil {
		logrus.Errorf("backup tar cf failed, ExecCommand err: %v", err)
		logrus.Infof("backup tar cf, errput: %s", errput)
		return
	}

	logrus.Infof("backup finished. clusterDir: %s", clusterDir)

	return
}

//add single node, for loop in AddNodes
func addSingleNode(nodeip string, sshUser string, privateKeyPath string, clusterDir string, slaveType string) {

	logrus.Infof("add node %s ...", nodeip)

	//scp -i $ssh_key $clusterDir/genconf/serve/dcos-install.tar $(sshuser)@$(nodeip):/tmp/dcos-install.tar
	commandStr := "scp -oConnectTimeout=10 -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null -oBatchMode=yes -oPasswordAuthentication=no -i " + privateKeyPath + " " +
		clusterDir + "genconf/serve/dcos-install.tar " + sshUser + "@" + nodeip + ":/tmp/dcos-install.tar"

	logrus.Infof("execute command: %s", commandStr)
	_, errput, err := common.ExecCommand(commandStr)
	if err != nil {
		logrus.Errorf("AddNodes, ExecCommand err: %v", err)
		logrus.Infof("add node %s failed failed, errput: %s", nodeip, errput)
		return
	}

	commandStr = "sudo mkdir -p /opt/dcos_install_tmp"
	_, errput, err = common.SshExecCmdWithKey(nodeip, "22", sshUser, privateKeyPath, commandStr)
	if err != nil {
		logrus.Errorf("AddNodes, ExecCommand err: %v", err)
		logrus.Infof("add node %s failed failed, errput: %s", nodeip, errput)
		return
	}

	commandStr = "sudo tar xf /tmp/dcos-install.tar -C /opt/dcos_install_tmp"
	_, errput, err = common.SshExecCmdWithKey(nodeip, "22", sshUser, privateKeyPath, commandStr)
	if err != nil {
		logrus.Errorf("AddNodes, ExecCommand err: %v", err)
		logrus.Infof("add node %s failed failed, errput: %s", nodeip, errput)
		return
	}

	commandStr = "sudo bash /opt/dcos_install_tmp/dcos_install.sh " + slaveType + " >> " + clusterDir + "opendcos_addnode.log"
	_, errput, err = common.SshExecCmdWithKey(nodeip, "22", sshUser, privateKeyPath, commandStr)
	if err != nil {
		logrus.Errorf("AddNodes, ExecCommand err: %v", err)
		logrus.Infof("add node %s failed failed, errput: %s", nodeip, errput)
		return
	}
	logrus.Infof("add node %s succeeded", nodeip)
	return
}

//delete single node
func deleteSingleNode(nodeip string, sshUser string, privateKeyPath string) {

	logrus.Infof("delete node %s ...", nodeip)

	commandStr := "sudo -i /opt/mesosphere/bin/pkgpanda uninstall"
	output, errput, err := common.SshExecCmdWithKey(nodeip, "22", sshUser, privateKeyPath, commandStr)
	if err != nil {
		logrus.Errorf("DeleteNode, ExecCommand err: %v", err)
		logrus.Infof("DeleteNode failed, errput: %s", errput)
		return
	}
	logrus.Infof("command output: %s", output)

	commandStr = "sudo rm -rf /opt/mesosphere /etc/mesosphere"
	_, errput, err = common.SshExecCmdWithKey(nodeip, "22", sshUser, privateKeyPath, commandStr)
	if err != nil {
		logrus.Errorf("DeleteNode, ExecCommand err: %v", err)
		logrus.Infof("DeleteNode failed, errput: %s", errput)
		return
	}
}

//TODO
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
	logrus.Infof("genConf, slaveStr:\n%s", strSlaves)
	logrus.Infof("genConf, masterStr:\n%s", strMasters)

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

	logrus.Debugf("genConf, config.yaml:\n", fileStr)
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

	logrus.Debugf("genIPDetect, ip-detect:\n%s", fileStr)
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

	logrus.Infof("start create file ssh_key... path: %s, privateKey:\n%s", path, privateKey)

	decodeKey, err := common.Base64Decode([]byte(privateKey))
	if err != nil {
		logrus.Errorf("decode private key error %v", err)
		return
	}

	logrus.Infof("the decoded privateKey:\n%s", decodeKey)

	err = ioutil.WriteFile(path+"ssh_key", []byte(decodeKey), 0600)
	if err != nil {
		logrus.Errorf("genSshKey failed, err is %v", err)
	}
	logrus.Infof("file ssh_key created.")
	return
}

//generate cluster directory with username & clusterName
func genClusterDir(username string, clusterName string) (clusterDir string) {

	clusterDir = BASE_PATH + username + "/" + clusterName + "/"
	return
}

//download dcos_generate_config.sh
func DownloadInstaller() (err error) {

	logrus.Infof("DownloadInstaller, downloading dcos_generate_config.sh. this may take some time...")

	//check if file exists
	exist, _ := common.PathExist("/opendcos/dcos_generate_config.sh")
	if exist {
		logrus.Infof("DownloadInstaller, dcos_generate_config.sh existed.")
		return
	}

	commandStr := "mkdir -p /opendcos/"
	logrus.Infof("execute command: %s", commandStr)
	_, errput, err := common.ExecCommand(commandStr)
	if err != nil {
		logrus.Errorf("DownloadInstaller, ExecCommand err: %v", err)
		logrus.Infof("DownloadInstaller failed, errput: %s", errput)
		return
	}

	commandStr = "curl -o /opendcos/dcos_generate_config.sh https://downloads.dcos.io/dcos/stable/dcos_generate_config.sh"
	logrus.Infof("execute command: %s", commandStr)
	output, errput, err := common.ExecCommand(commandStr)
	logrus.Infof(output)
	if err != nil {
		logrus.Errorf("DownloadInstaller, ExecCommand err: %v", err)
		logrus.Infof("DownloadInstaller failed, errput: %s", errput)
		return
	}
	return

}
