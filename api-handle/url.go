package api_handle

import "github.com/gin-gonic/gin"

func setupRoute(server *gin.Engine, webContext *WebProxy) {
	v1 := server.Group("/api/v1")
	{
		//----------------CCU-------------------
		control := v1.Group("/ctl")
		{
			master := control.Group("/master")
			{
				master.GET("", webContext.control)
			}
			worker := control.Group("/worker")
			{
				worker.GET("", webContext.executer)
			}
		}

	}
}


