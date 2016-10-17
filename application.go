package hipsterCache

import (
	"fmt"
	"net/http"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus"

	"hipster-cache/common"
	"hipster-cache/config"
	"hipster-cache/tcp"
)

type Application struct {
	config      *config.Config
	logger      common.ILogger
	consul      *consulapi.Client
	cacheServer *tcp.CacheServer
}

func NewApplication(config *config.Config, logger common.ILogger) *Application {
	return &Application{config: config, logger: logger}
}

func (this *Application) Init() error {
	fmt.Printf("\n startInit")
	err := this.initDiscovery()
	if err != nil {
		return err
	}

	err = this.initTCP()
	if err != nil {
		return err
	}
	this.initRouting()
	return nil
}

func (this *Application) initRouting() {
	// Handler for Prometheus
	http.Handle("/metrics", prometheus.Handler())
}

func (this *Application) Run() error {
	go http.ListenAndServe(fmt.Sprintf(":%d", this.config.MetricsPort), nil)
	this.cacheServer.Run()
	return nil
}

func (this *Application) registerService(catalog *consulapi.Catalog, id string, serviceName string, port int) error {
	service := &consulapi.AgentService{
		ID:      id,
		Service: serviceName,
		Port:    port,
	}

	reg := &consulapi.CatalogRegistration{
		Datacenter: "dc1",
		Node:       id,
		Address:    this.config.Address,
		Service:    service,
	}
	_, err := catalog.Register(reg, nil)
	return err
}

func (this *Application) initDiscovery() error {
	var err error
	config := consulapi.DefaultConfig()
	config.Address = this.config.ConsulAddress
	this.consul, err = consulapi.NewClient(config)
	if err != nil {
		return err
	}

	catalog := this.consul.Catalog()

	// Register for Applications
	err = this.registerService(catalog, "cache1", "hipster-cache", this.config.ServerPort)

	if err != nil {
		return err
	}

	// Register for Prometheus
	return this.registerService(catalog, "cache1-mertics", "hipster-cache-metrics", this.config.MetricsPort)
}

func (this *Application) initTCP() error {
	this.cacheServer = tcp.NewCacheServer(this.config.ServerPort, this.logger)
	return this.cacheServer.InitConnection()
	/*
		fmt.Printf("TCP IP")
		ips,err := net.LookupIP("consul6")
		fmt.Printf("\n Finish \n")
		if err != nil {
			this.logger.Errorf(`Error get IP:"%s"`, err)
			return
		}
		if len(ips) == 0 {
		  this.logger.Errorf(`Error get IP`)
		  return
		}
	*/
	//	fmt.Printf("Ips: %#v", ips)
	//	ip := ips[0]
	//	fmt.Printf(`Ip: %s"`, string(ip))
	//	listener, err := net.ListenTCP("tcp", tcpAddr)
}
