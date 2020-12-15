package main

import (
	"flag"
	"github.com/headend/share-module/configuration"
	"log"
	"time"
)

type WebProxy struct {
	ListenHost	string
	ListenPort	int16
	RequestTimeout	int
}

func main()  {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	confFilePtr := flag.String("c", "/opt/iptv/application.yml", "Configure file")
	flag.Parse()
	// load config
	var conf configuration.Conf
	if confFilePtr != nil {
		conf.ConfigureFile = *confFilePtr
	}
}


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