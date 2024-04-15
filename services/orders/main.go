package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/puzzles/services/basket/dal"
	"github.com/puzzles/services/basket/gen"
	"github.com/puzzles/services/basket/handler"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	viper.SetConfigFile(fmt.Sprintf("basket.%s.toml", env))
	viper.SetConfigType("toml")
	viper.SetEnvPrefix("puzzles")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db := dal.NewRedisDal(ctx,
		viper.GetString("redis.addr"),
		viper.GetString("redis.password"),
		viper.GetDuration("redis.duration"),
	)

	handler := handler.NewBasketHandler(
		db,
	)

	swagger, err := gen.GetSwagger()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.Use(middleware.OapiRequestValidator(swagger))
	gen.RegisterHandlers(router, handler)

	http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("port")), router)
}
