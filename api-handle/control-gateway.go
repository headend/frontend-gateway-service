package api_handle

import (
	"context"
	"fmt"
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

	ctlRequestHandle(w, ctx, static_config.StartWorker, true, 0)

}

func (w *WebProxy) stopWorker(ctx *gin.Context) {

	ctlRequestHandle(w, ctx, static_config.StopWorker, false, 0)
}


func (w *WebProxy) updateWorker(ctx *gin.Context) {

	ctlRequestHandle(w, ctx, static_config.UpdateWorker, true, 0)
}

func (w *WebProxy) activeMonitorSignal(ctx *gin.Context) {

	ctlRequestHandle(w, ctx, static_config.StartMonitorSignal, true, 0)
}

func (w *WebProxy) inactiveMonitorSignal(ctx *gin.Context) {

	ctlRequestHandle(w, ctx, static_config.StopMonitorSignal, false, 0)
}

func (w *WebProxy) activeMonitorVideo(ctx *gin.Context) {

	ctlRequestHandle(w, ctx, static_config.StartMonitorVideo, true, 0)
}

func (w *WebProxy) inactiveMonitorVideo(ctx *gin.Context) {

	ctlRequestHandle(w, ctx, static_config.StopMonitorVideo, false, 0)
}

func (w *WebProxy) activeMonitorAudio(ctx *gin.Context) {

	ctlRequestHandle(w, ctx, static_config.StartMonitorAudio, true, 0)
}

func (w *WebProxy) inactiveMonitorAudio(ctx *gin.Context) {

	ctlRequestHandle(w, ctx, static_config.StopMonitorAudio, false, 0)
}

func (w *WebProxy) updateThreadMonitor(ctx *gin.Context) {

	ctlRequestHandle(w, ctx, static_config.RunThread, false, 0)
}

func ctlRequestHandle(w *WebProxy, ctx *gin.Context, controlType int, isActive bool, thread_num int) {
	var requestData model.AgentCtlRequest
	err := ctx.BindJSON(&requestData)
	if err != nil {
		log.Println(err)
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
	if controlType != static_config.StartWorker {
		if !res.Agents[0].IsMonitor{
			// Phải bật monitor trước khi bật mấy cái khác
			ctx.String(401, "Enable worker first")
			return
		}
	}

	// Make Queue message for control request
	var queueFail bool
	switch controlType {
	case static_config.StartWorker:
		if res.Agents[0].SignalMonitor {
			queueFail = sendToQueue(w, static_config.StartMonitorSignal, requestData, isActive, thread_num)
		}
		if res.Agents[0].VideoMonitor {
			queueFail = sendToQueue(w, static_config.StartMonitorVideo, requestData, isActive, thread_num)
		}
	case static_config.StopWorker:
		if res.Agents[0].SignalMonitor {
			queueFail = sendToQueue(w, static_config.StopMonitorSignal, requestData,  false, thread_num)
		}
		if res.Agents[0].VideoMonitor {
			queueFail = sendToQueue(w, static_config.StopMonitorVideo, requestData, false, thread_num)
		}
	default:
		queueFail = sendToQueue(w, controlType, requestData, isActive, thread_num)
	}
	if queueFail {
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
		AgentId: requestData.AgentId,
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
	// update db
	if err5 := agentCTLUpdateDB(w, controlType, requestData, isActive, thread_num);err != nil {
		log.Println(err5)
		ctx.String(500, "Internal Server error")
		return
	}
	ctx.String(200, rspDataString)
	return
}


func sendToQueue(w *WebProxy, controlType int, requestData model.AgentCtlRequest, isActive bool, threadNum int) (fail bool) {
	dataSendQueue := agentctlpb.AgentCTLRequest{
		AgentId:     requestData.AgentId,
		ControlType: int64(controlType),
		TunnelData:  nil,
		RunThread: int64(threadNum),
	}
	//log.Printf("%#v", dataSendQueue)
	// call rpc
	c, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var rpcErr error
	_, rpcErr = (*w.agentctlclient).ControlAgent(c, &dataSendQueue)
	if rpcErr != nil {
		log.Println(rpcErr.Error())
		return true
	}
	return false
}

func agentCTLUpdateDB(w *WebProxy, controlType int, requestData model.AgentCtlRequest, isActive bool, threadNum int) (err error) {
	// call rpc to updte data
	c, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var rpcErr error
	switch controlType {
	case static_config.StartWorker:
		_, rpcErr = (*w.agentclient).UpdateActiveMonitor(c, &agentpb.AgentActiveMonitor{
			Id:        requestData.AgentId,
			IsMonitor: true,
		})
	case static_config.StopWorker:
		_, rpcErr = (*w.agentclient).UpdateActiveMonitor(c, &agentpb.AgentActiveMonitor{Id: requestData.AgentId, IsMonitor: false})
	case static_config.UpdateWorker:
		_, rpcErr = (*w.agentclient).UpdateRunthread(c, &agentpb.AgentUpdateMonitorRunThread{
			Id:        requestData.AgentId,
		})
	case static_config.StartMonitorSignal:
		_, rpcErr = (*w.agentclient).UpdateMonitorSignal(c, &agentpb.AgentActiveMonitorSignal{
			Id:            requestData.AgentId,
			SignalMonitor: true,
		})
	case static_config.StopMonitorSignal:
		_, rpcErr = (*w.agentclient).UpdateMonitorSignal(c, &agentpb.AgentActiveMonitorSignal{
			Id:            requestData.AgentId,
			SignalMonitor: false,
		})
	case static_config.StartMonitorVideo:
		_, rpcErr = (*w.agentclient).UpdateMonitorVideo(c, &agentpb.AgentActiveMonitorVideo{
			Id:           requestData.AgentId,
			VideoMonitor: true,
		})
	case static_config.StopMonitorVideo:
		_, rpcErr = (*w.agentclient).UpdateMonitorVideo(c, &agentpb.AgentActiveMonitorVideo{
			Id:           requestData.AgentId,
			VideoMonitor: false,
		})
	case static_config.RunThread:
		_, rpcErr = (*w.agentclient).UpdateRunthread(c, &agentpb.AgentUpdateMonitorRunThread{
			Id:        requestData.AgentId,
			RunThread: int64(requestData.RunThread),
		})
	default:
		resultErr := fmt.Errorf("Can not match control type %d", controlType)
		log.Println(resultErr)
		return resultErr
	}
	if rpcErr != nil {
		log.Println(rpcErr.Error())
		return rpcErr
	}
	return nil
}

