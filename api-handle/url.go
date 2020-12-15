package api_handle

import "github.com/gin-gonic/gin"

func setupRoute(server *gin.Engine, webContext *WebProxy) {
	v1 := server.Group("/api/v1")
	{
		//----------------CCU-------------------
		users := v1.Group("/getCCU")
		{
			users.GET("", webContext.getCCU)
		}

	}
}


