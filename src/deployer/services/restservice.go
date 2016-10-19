package services

import (
	"deployer/common/entity"
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
)

type Resource struct {
}

type RespStruct struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Err     string      `json:"err"`
}

type RespData struct {
}

var (
	DEPLOY_ERROR_PARSE_REQUESTBODY_FAILED string = "PARSE_REQUESTBODY_FAILED"
	DEPLOY_ERROR_CREATECLUSTER_FAILED     string = "CREATECLUSTER_FAILED"
)

func (r Resource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/v1/deploy").
		Consumes("*/*").
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/cluster").To(r.CreateClusterHandler).
		//docs
		Doc("create a cluster").
		Operation("CreateClusterHandler").
		Param(ws.BodyParameter("body", "entity.CreateRequest").DataType("string")))

	ws.Route(ws.POST("/nodes").To(r.AddNodesHandler).
		//docs
		Doc("add nodes").
		Operation("AddNodesHandler").
		Param(ws.BodyParameter("body", "entity.AddNodeRequest").DataType("string")))

	ws.Route(ws.DELETE("/cluster/{username}/{clustername}").To(r.DeleteClusterHandler).
		//docs
		Doc("delete a cluster").
		Operation("DeleteClusterHandler").
		Param(ws.PathParameter("username", "username").DataType("string")).
		Param(ws.PathParameter("clustername", "clustername").DataType("string")))

	ws.Route(ws.DELETE("/nodes/{username}/{clustername}/{nodeip}").To(r.DeleteNodeHandler).
		//docs
		Doc("delete a node").
		Operation("DeleteNodeHandler").
		Param(ws.PathParameter("username", "username").DataType("string")).
		Param(ws.PathParameter("clustername", "clustername").DataType("string")).
		Param(ws.PathParameter("nodeip", "ip of node").DataType("string")))

	container.Add(ws)
}

func (r Resource) CreateClusterHandler(request *restful.Request, response *restful.Response) {

	logrus.Infof("create a cluster...")

	createRequest := entity.CreateRequest{}

	err := json.NewDecoder(request.Request.Body).Decode(&createRequest)
	if err != nil {
		logrus.Errorf("CreateClusterHandler, convert body to request failed, error is %v", err)
		resp := RespStruct{Success: false, Err: DEPLOY_ERROR_PARSE_REQUESTBODY_FAILED}
		response.WriteEntity(resp)
		return
	}

	err = CreateCluster(createRequest)
	if err != nil {
		logrus.Errorf("CreateClusterHandler, CreateCluster failed, error is %v", err)
		resp := RespStruct{Success: false, Err: DEPLOY_ERROR_CREATECLUSTER_FAILED}
		response.WriteEntity(resp)
		return
	}

	respData := RespData{}
	resp := RespStruct{Success: true, Data: respData}
	response.WriteEntity(resp)
	return
}

func (r Resource) AddNodesHandler(request *restful.Request, response *restful.Response) {

	logrus.Infof("add nodes...")

	respData := RespData{}
	resp := RespStruct{Success: true, Data: respData}
	response.WriteEntity(resp)
}

func (r Resource) DeleteClusterHandler(request *restful.Request, response *restful.Response) {

	logrus.Infof("delete a cluster...")
	username := request.PathParameter("username")
	clustername := request.PathParameter("clustername")
	logrus.Infof(username, clustername)

	respData := RespData{}
	resp := RespStruct{Success: true, Data: respData}
	response.WriteEntity(resp)
}

func (r Resource) DeleteNodeHandler(request *restful.Request, response *restful.Response) {

	logrus.Infof("delete a node...")
	username := request.PathParameter("username")
	clustername := request.PathParameter("clustername")
	nodeip := request.PathParameter("nodeip")
	logrus.Infof(username, clustername, nodeip)

	respData := RespData{}
	resp := RespStruct{Success: true, Data: respData}
	response.WriteEntity(resp)
}
