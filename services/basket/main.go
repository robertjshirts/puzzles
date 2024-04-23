package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"

	"github.com/puzzles/services/basket/dal"
	"github.com/puzzles/services/basket/gen"
	"github.com/puzzles/services/basket/handler"
)

func main() {
	// Determine environment
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	// Load configuration
	viper.SetConfigFile(fmt.Sprintf("basket.%s.toml", env))
	viper.SetConfigType("toml")
	viper.SetEnvPrefix("puzzles")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	// Connect to redis
	fmt.Printf("Connecting to redis on: %s", viper.GetString("redis.host"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db := dal.NewRedisDal(ctx,
		viper.GetString("redis.host"),
		viper.GetString("redis.port"),
		viper.GetString("redis.password"),
		viper.GetDuration("redis.duration"),
	)

	// Create handler
	handler := handler.NewBasketHandler(
		db,
	)
	swagger, err := gen.GetSwagger()
	if err != nil {
		panic(err)
	}

	// Create router
	router := gin.Default()
	router.Use(middleware.OapiRequestValidator(swagger))
	gen.RegisterHandlers(router, handler)

	// Start server
	http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("port")), router)
}
