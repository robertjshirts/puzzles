package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	viper.SetConfigFile(fmt.Sprintf("orders.%s.toml", env))
	viper.SetConfigType("toml")
	viper.SetEnvPrefix("puzzles")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("port")), router)
}
