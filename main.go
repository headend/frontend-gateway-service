package main

import (
	"flag"
	"github.com/headend/share-module/configuration"
	"log"
	"github.com/headend/frontend-gateway-service/api-handle"
)

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
