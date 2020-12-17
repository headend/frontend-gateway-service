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
func (w *WebProxy) startWorker(ctx *gin.Context) {

	ctlRequestHandle(ctx)
}

func (w *WebProxy) stopWorker(ctx *gin.Context) {

	ctlRequestHandle(ctx)
}


func (w *WebProxy) updateWorker(ctx *gin.Context) {

	ctlRequestHandle(ctx)
}


func ctlRequestHandle(ctx *gin.Context) {
	var requestData model.AgentCtlRequest
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
	filee.Path = static_config.LogPath + "control_message"
	filee.WriteString(string(b))
	responseData := model.AgentCtlDataResponse{
		ControlId:   requestData.ControlId,
	}
	responseResult := model.AgentCtlResponse{
		ReturnCode:    1,
		ReturnData:    responseData,
		ReturnMessage: "Success",
		TunnelData:    requestData.TunnelData,
	}
	rspDataString, err := getCtlResponseData(responseResult)
	if err != nil {
		log.Println(err)
		ctx.String(500, "Internal Server error")
		return
	}
	ctx.String(200, rspDataString)
	return
}

func getCtlResponseData(responseResult model.AgentCtlResponse) (rspDataString string, err error) {
	b, err := json.Marshal(responseResult)
	if err != nil {
		return "", err
	}
	return string(b), nil
}