package api_handle

import (
	"context"
	"github.com/gin-gonic/gin"
	agentpb "github.com/headend/iptv-agent-service/proto"
	agentctlpb "github.com/headend/agent-control-service/proto"
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
func (w *WebProxy) startWorker(ctx *gin.Context) {

	ctlRequestHandle(w, ctx, static_config.StartWorker)
}

func (w *WebProxy) stopWorker(ctx *gin.Context) {

	ctlRequestHandle(w, ctx, static_config.StopWorker)
}


func (w *WebProxy) updateWorker(ctx *gin.Context) {

	ctlRequestHandle(w, ctx, static_config.UpdateWorker)
}


func ctlRequestHandle(w *WebProxy, ctx *gin.Context, controlType int) {
	var requestData model.AgentCtlRequest
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
	// Make Queue message for control request
	if sendToQueue(w, controlType, requestData) {
		responseResult := model.AgentCtlResponse{
			ReturnCode:    0,
			ReturnData:    model.AgentCtlDataResponse{},
			ReturnMessage: "RPC connection error",
			TunnelData:    nil,
		}
		rspString, _ := responseResult.GetJsonString()
		ctx.String(200, rspString)
		return
	}
	responseData := model.AgentCtlDataResponse{
		ControlId: requestData.ControlId,
	}
	responseResult := model.AgentCtlResponse{
		ReturnCode:    1,
		ReturnData:    responseData,
		ReturnMessage: "Success",
		TunnelData:    requestData.TunnelData,
	}
	rspDataString, err := responseResult.GetJsonString()
	if err != nil {
		log.Println(err)
		ctx.String(500, "Internal Server error")
		return
	}
	ctx.String(200, rspDataString)
	return
}


func sendToQueue(w *WebProxy, controlType int, requestData model.AgentCtlRequest) (fail bool) {
	dataSendQueue := agentctlpb.AgentCTLRequest{
		AgentId:     requestData.AgentId,
		ControlType: int64(controlType),
		TunnelData:  nil,
	}
	// call rpc
	c, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var rpcErr error
	switch controlType {
	case static_config.StartWorker:
		_, rpcErr = (*w.agentctlclient).STARTWorker(c, &dataSendQueue)
	case static_config.StopWorker:
		_, rpcErr = (*w.agentctlclient).STOPWorker(c, &dataSendQueue)
	case static_config.UpdateWorker:
		_, rpcErr = (*w.agentctlclient).STOPWorker(c, &dataSendQueue)
	default:
		log.Println("Can not match control type")
		return true
	}
	if rpcErr != nil {
		log.Println(rpcErr.Error())
		return true
	}
	return true
}
