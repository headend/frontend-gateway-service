package api_handle

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"github.com/headend/share-module/configuration"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	agentpb "github.com/headend/iptv-agent-service/proto"
	agentctlpb "github.com/headend/agent-control-service/proto"
	agentexepb "github.com/headend/agent-executer-service/proto"
)

type WebProxy struct {
	Conf *configuration.Conf
	agentclient	*agentpb.AgentServiceClient
	agentctlclient	*agentctlpb.AgentCTLServiceClient
	agentexeclient	*agentexepb.AgentEXEServiceClient
}

func StartAgentGatewayService(config *configuration.Conf)  {
	//connect agent services
	agentConn := initializeClient(config.RPC.Agent.Gateway, config.RPC.Agent.Port)
	defer agentConn.Close()
	agentClient := agentpb.NewAgentServiceClient(agentConn)

	//connect agent coltrol services
	agentCtlConn := initializeClient(config.RPC.Agentctl.Gateway, config.RPC.Agentctl.Port)
	defer agentCtlConn.Close()
	agentCtlClient := agentctlpb.NewAgentCTLServiceClient(agentCtlConn)

	//connect agent executer services
	agentExeConn := initializeClient(config.RPC.Agentrunner.Gateway, config.RPC.Agentrunner.Port)
	defer agentExeConn.Close()
	agentExeClient := agentexepb.NewAgentEXEServiceClient(agentExeConn)

	webContext := WebProxy{
		Conf: config,
		agentclient: &agentClient,
		agentctlclient: &agentCtlClient,
		agentexeclient: &agentExeClient,
	}
	server := initializeServer(config.Server.RequestTimeout)
	setupRoute(server, &webContext)
	listenAdd := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	swgrUri := fmt.Sprintf("%s/swagger/doc.json", listenAdd)
	swgUrl := ginSwagger.URL(swgrUri)
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swgUrl))
	log.Print("begin run http server...")
	log.Printf("Visit document page: %s/swagger/", listenAdd)
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

func initializeClient(host string, port uint16) *grpc.ClientConn {
	connectAddr := fmt.Sprintf("%s:%d", host, port)
	println(connectAddr)
	conn, err := grpc.Dial(connectAddr,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(5*1024*1024)),
		grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect server: %v", err)
	}
	return conn
}




