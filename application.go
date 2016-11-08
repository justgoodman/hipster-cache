package hipsterCache

import (
	"fmt"
	"net/http"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus"

	"hipster-cache/common"
	"hipster-cache/config"
	"hipster-cache/hash_table"
	"hipster-cache/tcp"
)

type Application struct {
	config      *config.Config
	logger      common.ILogger
	consul      *consulapi.Client
	cacheServer *tcp.CacheServer
	hashTable   *hash_table.HashTable
}

func NewApplication(config *config.Config, logger common.ILogger) *Application {
	return &Application{config: config, logger: logger}
}

func (a *Application) Init() error {
	fmt.Printf("\n startInit")
	err := a.initDiscovery()
	if err != nil {
		return err
	}
	a.initHashTable()

	err = a.initTCP()
	if err != nil {
		return err
	}
	a.initRouting()
	return nil
}

func (a *Application) initRouting() {
	// Handler for Prometheus
	http.Handle("/metrics", prometheus.Handler())
}

func (a *Application) initHashTable() {
	a.hashTable = hash_table.NewHashTable(a.config.InitCapacity, a.config.MaxLenghtKey, a.config.MaxBytesSize)
	a.hashTable.InitMetrics()
}

func (a *Application) Run() error {
	go http.ListenAndServe(fmt.Sprintf(":%d", a.config.MetricsPort), nil)
	a.cacheServer.Run()
	return nil
}

func (a *Application) registerService(catalog *consulapi.Catalog, serviceName, nodeName string, port int) error {
	service := &consulapi.AgentService{
		ID:      serviceName,
		Service: serviceName,
		Port:    port,
	}

	reg := &consulapi.CatalogRegistration{
		Datacenter: "dc1",
		Node:       nodeName,
		Address:    a.config.Address,
		Service:    service,
		TaggedAddresses: map[string]string{
			"lan": a.config.Address,
			"wan": a.config.WANAddress,
		},
	}
	fmt.Printf("\n Registered '%#v' \n", reg)
	_, err := catalog.Register(reg, nil)
	return err
}

func (a *Application) registerHealthCheck(agent *consulapi.Agent, address string, port int) error {
	reg := &consulapi.AgentCheckRegistration{
		ID:   fmt.Sprintf("HealthCheckServer_%s", address),
		Name: fmt.Sprintf("Health Check TCP for node: %s", address),
	}
	reg.TCP = fmt.Sprintf("%s:%d", address, port)
	reg.Interval = "30s"
	return agent.CheckRegister(reg)
}

func (a *Application) initDiscovery() error {
	var err error
	config := consulapi.DefaultConfig()
	config.Address = a.config.ConsulAddress
	a.consul, err = consulapi.NewClient(config)
	if err != nil {
		return err
	}

	catalog := a.consul.Catalog()

	// Register for Applications
	err = a.registerService(catalog, "hipster-cache", a.config.Address, a.config.ServerPort)

	if err != nil {
		return err
	}

	agent := a.consul.Agent()
	// Register heath check
	err = a.registerHealthCheck(agent, a.config.Address, a.config.ServerPort)

	if err != nil {
		return err
	}

	// Register for Prometheus
	return a.registerService(catalog, "hipster-cache-metrics", a.config.Address, a.config.MetricsPort)
}

func (a *Application) initTCP() error {
	a.cacheServer = tcp.NewCacheServer(a.hashTable, a.logger, a.config.ServerPort)
	return a.cacheServer.InitConnection()
	/*
		fmt.Printf("TCP IP")
		ips,err := net.LookupIP("consul6")
		fmt.Printf("\n Finish \n")
		if err != nil {
			a.logger.Errorf(`Error get IP:"%s"`, err)
			return
		}
		if len(ips) == 0 {
		  a.logger.Errorf(`Error get IP`)
		  return
		}
	*/
	//	fmt.Printf("Ips: %#v", ips)
	//	ip := ips[0]
	//	fmt.Printf(`Ip: %s"`, string(ip))
	//	listener, err := net.ListenTCP("tcp", tcpAddr)
}
