package api_handle

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/headend/share-module/model"
	"github.com/headend/share-module/file-and-directory"
	"github.com/headend/share-module/configuration/static-config"
	"log"
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

	exeRequestHandle(ctx)
}


func exeRequestHandle(ctx *gin.Context) {
	var requestData model.AgentExeRequest
	err := ctx.BindJSON(&requestData)
	if err != nil {
		ctx.String(400, "invalid param")
		return
	}
	b, err := json.Marshal(requestData)
	if err != nil {
		ctx.String(400, "invalid param")
		return
	}
	var filee file_and_directory.MyFile
	filee.Path = static_config.LogPath + "executer_message"
	filee.WriteString(string(b))
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
