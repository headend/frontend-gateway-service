package api_handle

import (
	"fmt"
	"log"
	"github.com/headend/share-module/configuration"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type WebProxy struct {
	ListenHost	string
	ListenPort	uint16
	RequestTimeout	int
}

func StartAgentGatewayService(config *configuration.Conf)  {
	webContext := WebProxy{
		ListenHost:     config.Server.Host,
		ListenPort:     config.Server.Port,
		RequestTimeout: config.Server.RequestTimeout,
	}
	server := initializeServer(config.Server.RequestTimeout)
	setupRoute(server, &webContext)
	server.GET("/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Print("begin run http server...")
	listenAdd := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Printf("Visit document page: %s/doc/", listenAdd)
	_ = server.Run(listenAdd)

}

func initializeServer(RequestTimeout int) *gin.Engine {
	server := gin.New()
	gin.SetMode(gin.ReleaseMode)
	server.Use(gin.Logger())
	server.Use(gin.Recovery())

	// CORS for https://foo.com and https://github.com origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 30 seconds
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           time.Duration(RequestTimeout) * time.Second,
	}))
	return server
}
