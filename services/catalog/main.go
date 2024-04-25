package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"

	"github.com/puzzles/services/catalog/dal"
	"github.com/puzzles/services/catalog/gen"
	"github.com/puzzles/services/catalog/handler"
)

func main() {
	// Determine environment
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	// Load configuration
	viper.SetConfigFile(fmt.Sprintf("catalog.%s.toml", env))
	viper.SetConfigType("toml")
	viper.SetEnvPrefix("puzzles")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(viper.GetDuration("startup_delay"))

	// Connect to consul
	if env == "prod" {
		id := RegisterService()
		defer DeregisterService(id)
	}

	// Connect to mongo
	fmt.Printf("conn: %s\n", viper.GetString("mongo.uri"))
	db := dal.NewMongoDAL(
		viper.GetString("mongo.uri"),
		viper.GetString("mongo.db"),
	)

	// Create handler
	handler := handler.NewCatalogHandler(db)
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
