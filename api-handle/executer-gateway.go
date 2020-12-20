package api_handle

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	agentpb "github.com/headend/iptv-agent-service/proto"
	agentexepb "github.com/headend/agent-executer-service/proto"
	"github.com/headend/share-module/configuration/static-config"
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

	exeRequestHandle(w,ctx, static_config.UrgentTask)
}


func exeRequestHandle(w *WebProxy, ctx *gin.Context, exeType int) {
	var requestData model.AgentExeRequest
	err := ctx.BindJSON(&requestData)
	if err != nil {
		ctx.String(400, "invalid param")
		return
	}
	// Check Agent exists
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := (*w.agentclient).Get(c, &agentpb.AgentFilter{Id: requestData.AgentId})
	if err != nil {
		log.Println(err)
		ctx.String(500, "Internal server error")
		return
	}
	if len(res.Agents) == 0 {
		ctx.String(404, "Agent not found")
		return
	}
	// agent exist here

	// Send request to message queue
	if senMessageToQueue(w, exeType, requestData) {
		return
	}
	// Make data response
	responseData := model.AgentExeDataResponse{
		ExeType:   exeType,
		ExeId:     requestData.ExeId,
		ProfileId: requestData.ProfileId,
		AgentId:   requestData.AgentId,
	}
	responseResult := model.AgentExeResponse{
		ReturnCode:    1,
		ReturnData:    responseData,
		ReturnMessage: "Success",
		TunnelData:    requestData.TunnelData,
	}
	rspDataString, err2 := responseResult.GetJsonString()
	if err != nil {
		log.Println(err2)
		ctx.String(500, "Internal Server error")
		return
	}
	ctx.String(200, rspDataString)
	return
}

func senMessageToQueue(w *WebProxy, exeType int, requestData model.AgentExeRequest) (fail bool) {
	exeRequestData := agentexepb.AgentEXERequest{
		AgentId:    requestData.AgentId,
		ProfileId:  requestData.ProfileId,
		ExeId:      int64(requestData.ExeId),
		ExeType:    int64(exeType),
		TunnelData: nil,
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var err error
	switch exeType {
	case static_config.UrgentTask:
		_, err = (*w.agentexeclient).RunUrgentTask(c, &exeRequestData)
	case static_config.CommandShell:
		log.Println("Wait for support")
	default:
		log.Println("Unsupport executer type: ", exeType)
		err = fmt.Errorf("Unsupport executer type: %d", exeType)
	}

	if err != nil{
		log.Println(err)
		return true
	}
	return false
}



