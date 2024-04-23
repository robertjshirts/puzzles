package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	capi "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

func RegisterService() string {
	consulConfig := capi.DefaultConfig()
	consulConfig.Address = viper.GetString("consul.address")
	consulClient, err := capi.NewClient(consulConfig)
	if err != nil {
		log.Fatal(err)
	}

	tags := []string{
		"traefik.enable=true",
		fmt.Sprintf("traefik.http.routers.%s.rule=PathPrefix(`%s`)", viper.GetString("service.name"), viper.GetString("service.path")),
	}

	id := fmt.Sprintf("%s-%s", viper.GetString("service.name"), uuid.New().String())
	serviceReg := &capi.AgentServiceRegistration{
		ID:      id,
		Tags:    tags,
		Name:    viper.GetString("service.name"),
		Port:    viper.GetInt("port"),
		Address: viper.GetString("service.address"),
		Check: &capi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", viper.GetString("service.address"), viper.GetInt("port")),
			Interval: "10s",
			Timeout:  "3s",
		},
	}

	err = consulClient.Agent().ServiceRegister(serviceReg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Service %s with id %s registered in consul", viper.GetString("service.name"), id)

	return id
}

func DeregisterService(id string) {
	consulConfig := capi.DefaultConfig()
	consulConfig.Address = viper.GetString("consul.address")

	consulClient, err := capi.NewClient(consulConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = consulClient.Agent().ServiceDeregister(id)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Service with id %s deregistered from consul", id)
}
