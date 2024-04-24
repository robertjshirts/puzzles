package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	middleware "github.com/oapi-codegen/gin-middleware"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"

	"github.com/puzzles/services/order/dal"
	"github.com/puzzles/services/order/gen"
	"github.com/puzzles/services/order/handler"
)

func main() {
	// Determine environment
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	// Load configuration
	viper.SetConfigFile(fmt.Sprintf("order.%s.toml", env))
	viper.SetConfigType("toml")
	viper.SetEnvPrefix("puzzles")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	// Register service
	if env == "prod" {
		id := RegisterService()
		defer DeregisterService(id)
	}

	// Connect to postgres
	fmt.Printf("Connecting to postgres with: host=%s port=%d user=%s dbname=%s sslmode=disable\n", viper.GetString("postgres.host"), viper.GetInt("postgres.port"), viper.GetString("postgres.user"), viper.GetString("postgres.dbname"))
	db, err := dal.NewSQLDal(
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.user"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.dbname"),
		viper.GetDuration("postgres.timeout"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create handler
	handler := handler.NewOrderHandler(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.pass"),
		viper.GetString("rabbitmq.host"),
		viper.GetInt("rabbitmq.port"),
		viper.GetString("rabbitmq.queue"),
		db)
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
