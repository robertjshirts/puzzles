package main

import (
	"log"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"

	"github.com/puzzles/services/catalog/api"
	"github.com/puzzles/services/catalog/dal"
	"github.com/puzzles/services/catalog/gen"
)

func main() {
	viper.SetConfigName("catalog")
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetEnvPrefix("puzzles")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	db := dal.NewMongoDAL(
		viper.GetString("mongo.uri"),
		viper.GetString("mongo.db"),
	)

	handler := api.NewCatalogHandler(db)

	swagger, err := gen.GetSwagger()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.Use(middleware.OapiRequestValidator(swagger))
	gen.RegisterHandlers(router, handler)

	router.Run()
}
