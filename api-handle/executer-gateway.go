package api_handle

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	agentpb "github.com/headend/iptv-agent-service/proto"
	"github.com/headend/share-module/configuration/static-config"
	"github.com/headend/share-module/file-and-directory"
	"github.com/headend/share-module/model"
	"log"
	"time"
)


//
// @Summary Control agentd
// @Description Control agentd
// @Accept  json
// @Produce  json
// @param model.AgentCtlRequest body model.AgentCtlRequest true "Input params"
// @Success 200 {object} model.AgentCtlResponse    "True"
// @Failure 400 {string} string "Invalid param!!"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/ctl/master/ [post]
// @BasePath /v1
func (w *WebProxy) runUrgentTask(ctx *gin.Context) {

	exeRequestHandle(w,ctx)
}


func exeRequestHandle(w *WebProxy, ctx *gin.Context) {
	var requestData model.AgentExeRequest
	err := ctx.BindJSON(&requestData)
	if err != nil {
		ctx.String(400, "invalid param")
		return
	}
	// Check Agent exists
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := (*w.agentclient).Get(c, &agentpb.AgentFilter{IpControl: requestData.ListAgentToRun[0]})
	if err != nil {
		log.Println(err)
		ctx.String(500, "Internal server error")
		return
	}
	if len(res.Agents) == 0 {
		ctx.String(404, "Agent not found")
		return
	}

	jsonString, err := requestData.GetJsonString()
	if err != nil {
		log.Println(err)
		ctx.String(400, "Internal server error")
		return
	}
	var filee file_and_directory.MyFile
	filee.Path = static_config.LogPath + "executer_message"
	filee.WriteString(jsonString)
	responseData := model.AgentExeDataResponse{
		ExeType: requestData.ExeType,
		ExeId:   requestData.ExeId,
		ListAgentToRun: requestData.ListAgentToRun,
		ListProfileId: nil,
	}
	responseResult := model.AgentExeResponse{
		ReturnCode:    1,
		ReturnData:    responseData,
		ReturnMessage: "Success",
		TunnelData:    requestData.TunnelData,
	}
	rspDataString, err := getExeResponseData(responseResult)
	if err != nil {
		log.Println(err)
		ctx.String(500, "Internal Server error")
		return
	}
	ctx.String(200, rspDataString)
	return
}

func getExeResponseData(responseResult model.AgentExeResponse) (rspDataString string, err error) {
	b, err := json.Marshal(responseResult)
	if err != nil {
		return "", err
	}
	return string(b), nil
}


