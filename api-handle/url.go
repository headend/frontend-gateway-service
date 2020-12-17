package api_handle

import "github.com/gin-gonic/gin"

func setupRoute(server *gin.Engine, webContext *WebProxy) {
	v1 := server.Group("/api/v1")
	{
		//----------------CCU-------------------
		control := v1.Group("/ctl/:agent_ip")
		{
			master := control.Group("/master")
			{
				master.POST("/start-worker", webContext.startWorker)
				master.POST("/stop-worker", webContext.stopWorker)
				master.POST("/update-worker", webContext.updateWorker)
				monitor := master.Group("/monitor")
				{
					monitorSignal := monitor.Group("/signal")
					{
						monitorSignal.POST("/enable", webContext.startWorker)
						monitorSignal.POST("/disable", webContext.startWorker)
					}
					monitorVideo := monitor.Group("/video")
					{
						monitorVideo.POST("/enable", webContext.startWorker)
						monitorVideo.POST("/disable", webContext.startWorker)
					}
					monitorAudio := monitor.Group("/audio")
					{
						monitorAudio.POST("/enable", webContext.startWorker)
						monitorAudio.POST("/disable", webContext.startWorker)
					}
				}
			}
			worker := v1.Group("/exe/:agent_ip")
			{
				exeTask := worker.Group("/task")
				{
					exeTask.POST("/urgent", webContext.runUrgentTask)
				}
			}
		}

	}
}


