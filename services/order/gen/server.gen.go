// Package gen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package gen

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// Defines values for OrderStatus.
const (
	Cancelled OrderStatus = "cancelled"
	Placed    OrderStatus = "placed"
)

// Defines values for PuzzleType.
const (
	Clock    PuzzleType = "clock"
	Megaminx PuzzleType = "megaminx"
	N2x2     PuzzleType = "2x2"
	N3x3     PuzzleType = "3x3"
	N4x4     PuzzleType = "4x4"
	N5x5     PuzzleType = "5x5"
	N6x6     PuzzleType = "6x6"
	N7x7     PuzzleType = "7x7"
	N8x8     PuzzleType = "8x8+"
	Other    PuzzleType = "other"
	Pyraminx PuzzleType = "pyraminx"
	Skewb    PuzzleType = "skewb"
	Square1  PuzzleType = "square-1"
)

// Error defines model for Error.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewOrderInfo defines model for NewOrderInfo.
type NewOrderInfo struct {
	Items []Puzzle `json:"items"`

	// Name The name of the customer
	Name string `json:"name"`

	// PaymentInfo The payment information
	PaymentInfo PaymentInfo `json:"paymentInfo"`

	// ShippingInfo The shipping information
	ShippingInfo ShippingInfo `json:"shippingInfo"`
}

// OrderInfo defines model for OrderInfo.
type OrderInfo struct {
	// Id The ID of the order
	Id    string   `json:"id"`
	Items []Puzzle `json:"items"`

	// Name The name of the customer
	Name string `json:"name"`

	// PaymentInfo The payment information
	PaymentInfo PaymentInfo `json:"paymentInfo"`

	// ShippingInfo The shipping information
	ShippingInfo ShippingInfo `json:"shippingInfo"`

	// Status The status of the order
	Status OrderStatus `json:"status"`
}

// OrderStatus The status of the order
type OrderStatus string

// PaymentInfo The payment information
type PaymentInfo struct {
	// AreaCode The area code of the card owner
	AreaCode string `json:"areaCode"`

	// CardNumber The card number
	CardNumber string `json:"cardNumber"`

	// Cvv The CVV
	Cvv string `json:"cvv"`

	// Expiration The expiration date
	Expiration string `json:"expiration"`
}

// Puzzle A puzzle
type Puzzle struct {
	// Description A description of the puzzle
	Description string `json:"description"`

	// Id The ID of the puzzle
	Id string `json:"id"`

	// Name The name of the puzzle
	Name string `json:"name"`

	// Price The price of the puzzle
	Price float64 `json:"price"`

	// Type The type of puzzle
	Type PuzzleType `json:"type"`
}

// PuzzleType The type of puzzle
type PuzzleType string

// ShippingInfo The shipping information
type ShippingInfo struct {
	// Address The address of the recipient
	Address string `json:"address"`

	// AreaCode The area code of the recipient
	AreaCode string `json:"areaCode"`

	// City The city of the recipient
	City string `json:"city"`

	// Country The country of the recipient
	Country string `json:"country"`

	// Name The name of the recipient
	Name string `json:"name"`

	// State The state of the recipient
	State string `json:"state"`
}

// CreateNewOrderJSONRequestBody defines body for CreateNewOrder for application/json ContentType.
type CreateNewOrderJSONRequestBody = NewOrderInfo

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Check the health of the service
	// (GET /health)
	CheckHealth(c *gin.Context)
	// Place a new order
	// (POST /order)
	CreateNewOrder(c *gin.Context)
	// Delete an order by ID
	// (DELETE /order/{id})
	DeleteOrder(c *gin.Context, id string)
	// Get order details by ID
	// (GET /order/{id})
	GetOrder(c *gin.Context, id string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// CheckHealth operation middleware
func (siw *ServerInterfaceWrapper) CheckHealth(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CheckHealth(c)
}

// CreateNewOrder operation middleware
func (siw *ServerInterfaceWrapper) CreateNewOrder(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CreateNewOrder(c)
}

// DeleteOrder operation middleware
func (siw *ServerInterfaceWrapper) DeleteOrder(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.DeleteOrder(c, id)
}

// GetOrder operation middleware
func (siw *ServerInterfaceWrapper) GetOrder(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetOrder(c, id)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/health", wrapper.CheckHealth)
	router.POST(options.BaseURL+"/order", wrapper.CreateNewOrder)
	router.DELETE(options.BaseURL+"/order/:id", wrapper.DeleteOrder)
	router.GET(options.BaseURL+"/order/:id", wrapper.GetOrder)
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xWS3PbNhD+K5htb6Ut+hEno1trd1JfEs/Yk0vGB5hcSYhJAAFAm4qH/72zACiKD1nq",
	"NO34xAcWi93v2/2wL5CpUiuJ0lmYv4DNVlhy//qnMcrQizZKo3EC/e9M5UhPt9YIcxDS4RINNAmUaC1f",
	"bi9aZ4RcQtMkYPB7JQzmMP8aXHT290lrrx6+YebI1yd8/mxyNNdyocYxCIdl/+VXgwuYwy+zLp1ZzGV2",
	"U/34USB5jcdwY/iaviUvfbg52swI7YSSMIe7FTJaYWrB3ApZVlmnSjSQDPNKQPN1idK1Yb4axpZpk4Bd",
	"Ca2FXB6y9XbbdohmwCAmM/DbD3AK6B7KvCg+L2D+9fVoetw0yYicfBrS66sWUEW7p9C0jrtqL5/+8Ntg",
	"2jQUQA+OHDaO7pv7ZCISfz4TcqFMyf3/XcDcbgIaewlnDHNCWZUUhi54hhRKxmWGRYH5Fvxdxjf9+hmf",
	"EvkbRNtHnBvkl7Etxx5olVHLbeqZm5ypZznNAa1+qsoHNNPu/G4ZDKa2Pz1N77v88mXKHmstTEhrclu3",
	"znLucOxiqC1d+D3nIbKkw2qqGaJSjAL5nemwMkS+ZzbetfXdgr9xNEJif9/s3nuYkO3er43IdjjwSyMP",
	"Ja9FSZV+kqZpmkApZPim96pwQhdIUpIepyeb8zZVUx8t1VH8uSgUdxfnnTwfIuZ3ZDkSwrxTQe+q3/tt",
	"lruJv4vnj0GgDYTBJv+2zU/rU0jgrD6DBM7rc0jgXf0OErioLyCB9/V7SOBD/eE3On1teClk7W++Zftq",
	"H/H5gZ7fK27w6IQKtVDZIySg3ArNpGrcDu6OCXGKFq/rRp4btDvkLS62zBvMhBYo3VT5/EMBetVXJtx6",
	"h/IItz7MhaqkM7u8hMWDHB3WVq+6oFsCd18ghzgZlHms8Ja8LfQjeO2hHRDjkienIlaPE45ELzYBu+TG",
	"sVs0T9QsCTyhsSHm9Dg9PqGclEbJtYA5nB2nx2d+ynArX0ezFfLCreh1iY4eVHC+/K5zmMPlCrPHv4IN",
	"pWW1kjZU42ma7sApxMKEZcG7n91yXPCqcGEklY6Qo5LWuhCZP2/2zQZVDtKxT1jCuOuBGci4ZEhrTGVZ",
	"ZYgEsrFVWXKqsZCS5zBE1zIaw/bWszAe0CCr7BQsBrnDdrCCQDha94fK1z8tw/7cNk6UsPbDJHOK2p8e",
	"cV7wIjlg6+SnBbY3KoqAPXPLMo9T/qb4v6FJj3Em8TlOgR3jsxeRN6GqCww60Cf+yv9vWdfc8BIdGuuH",
	"8L0zNDEUPVM7w9w3YnsLzsON2GmHMxUmW3AMdeZ+xPH5dEeGw4mQcPrbIiRgyriMcT6s2fUVRTipSB/R",
	"/Rv4DToj8Om/IiD9/5os1u4bIvIjugh0jo6LwrZUNk3zdwAAAP//uW0B/rkQAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
