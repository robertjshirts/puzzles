package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	capi "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

func RegisterService() string {
	// Get consul client
	consulConfig := capi.DefaultConfig()
	consulConfig.Address = viper.GetString("consul.address")
	consulClient, err := capi.NewClient(consulConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Generate tags and id, and get hostname
	id := fmt.Sprintf("%s-%s", viper.GetString("service.name"), uuid.New().String())
	tags := []string{
		"traefik.enable=true",
		fmt.Sprintf("traefik.http.routers.%s.rule=PathPrefix(`%s`)", viper.GetString("service.name"), viper.GetString("service.path")),
	}
	host, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	// Register service with tags, id, and other configurations
	serviceReg := &capi.AgentServiceRegistration{
		ID:      id,
		Tags:    tags,
		Name:    viper.GetString("service.name"),
		Port:    viper.GetInt("port"),
		Address: host,
		Check: &capi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", host, viper.GetInt("port")),
			Interval: "10s",
			Timeout:  "3s",
		},
	}
	err = consulClient.Agent().ServiceRegister(serviceReg)
	if err != nil {
		log.Fatal(err)
	}

	// Return id
	log.Printf("Service %s with id %s registered in consul", viper.GetString("service.name"), id)
	return id
}

func DeregisterService(id string) {
	// Get consul client
	consulConfig := capi.DefaultConfig()
	consulConfig.Address = viper.GetString("consul.address")
	consulClient, err := capi.NewClient(consulConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Deregister service
	err = consulClient.Agent().ServiceDeregister(id)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Service with id %s deregistered from consul", id)
}
