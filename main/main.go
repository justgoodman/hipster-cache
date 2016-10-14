package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/juju/loggo"
	app "hipster-cache"
	"hipster-cache/config"
)

var addr = flag.String("listen-address", ":4001", "The address to listen on for HTTP requests.")

func main() {
	flag.Parse()
	config := config.NewConfig()

	logger := loggo.GetLogger("")

	err := config.LoadFile("etc/application.json")
	if err != nil {
		logger.Criticalf("Error reading configuration file: '%s'", err.Error())
		os.Exit(1)
	}
	logger.Errorf("Test Error")
	application := app.NewApplication(config, logger)
	fmt.Printf("#%v", application)
	err = application.Init()
	if err != nil {
		logger.Criticalf("error initialization application: '%s'", err.Error())
		os.Exit(1)
	}
	application.Run()
	//	http.Handle("/metrics", prometheus.Handler())
	//http.ListenAndServe(*addr, nil)
}
