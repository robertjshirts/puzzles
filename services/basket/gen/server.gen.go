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

// Basket A basket
type Basket struct {
	// Id The ID of the basket
	Id string `json:"id"`

	// Items The items in the basket
	Items []Puzzle `json:"items"`
}

// Error An error response
type Error struct {
	// Code A code representing the error
	Code int32 `json:"code"`

	// Message A message describing the error
	Message string `json:"message"`
}

// Puzzle defines model for Puzzle.
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

// CreateNewBasketJSONBody defines parameters for CreateNewBasket.
type CreateNewBasketJSONBody = []Puzzle

// AddItemToBasketJSONBody defines parameters for AddItemToBasket.
type AddItemToBasketJSONBody = []Puzzle

// CreateNewBasketJSONRequestBody defines body for CreateNewBasket for application/json ContentType.
type CreateNewBasketJSONRequestBody = CreateNewBasketJSONBody

// AddItemToBasketJSONRequestBody defines body for AddItemToBasket for application/json ContentType.
type AddItemToBasketJSONRequestBody = AddItemToBasketJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create a new basket
	// (POST /basket)
	CreateNewBasket(c *gin.Context)
	// Delete a basket by ID
	// (DELETE /basket/{id})
	DeleteBasket(c *gin.Context, id string)
	// Get a basket by ID
	// (GET /basket/{id})
	GetBasket(c *gin.Context, id string)
	// Add items to a basket
	// (POST /basket/{id})
	AddItemToBasket(c *gin.Context, id string)
	// Remove a puzzle from a basket
	// (DELETE /basket/{id}/{puzzleId})
	RemovePuzzleFromBasket(c *gin.Context, id string, puzzleId string)
	// Check the health of the service
	// (GET /health)
	CheckHealth(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// CreateNewBasket operation middleware
func (siw *ServerInterfaceWrapper) CreateNewBasket(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CreateNewBasket(c)
}

// DeleteBasket operation middleware
func (siw *ServerInterfaceWrapper) DeleteBasket(c *gin.Context) {

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

	siw.Handler.DeleteBasket(c, id)
}

// GetBasket operation middleware
func (siw *ServerInterfaceWrapper) GetBasket(c *gin.Context) {

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

	siw.Handler.GetBasket(c, id)
}

// AddItemToBasket operation middleware
func (siw *ServerInterfaceWrapper) AddItemToBasket(c *gin.Context) {

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

	siw.Handler.AddItemToBasket(c, id)
}

// RemovePuzzleFromBasket operation middleware
func (siw *ServerInterfaceWrapper) RemovePuzzleFromBasket(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "puzzleId" -------------
	var puzzleId string

	err = runtime.BindStyledParameterWithOptions("simple", "puzzleId", c.Param("puzzleId"), &puzzleId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter puzzleId: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.RemovePuzzleFromBasket(c, id, puzzleId)
}

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

	router.POST(options.BaseURL+"/basket", wrapper.CreateNewBasket)
	router.DELETE(options.BaseURL+"/basket/:id", wrapper.DeleteBasket)
	router.GET(options.BaseURL+"/basket/:id", wrapper.GetBasket)
	router.POST(options.BaseURL+"/basket/:id", wrapper.AddItemToBasket)
	router.DELETE(options.BaseURL+"/basket/:id/:puzzleId", wrapper.RemovePuzzleFromBasket)
	router.GET(options.BaseURL+"/health", wrapper.CheckHealth)
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xWTXPbNhD9K5htb4UlynacDG923Ka6tJ3WN08OELmSEBMAA4A2FY3+e2cBkpJFKnFr",
	"u6PO5MQPYBdv97195Boyo0qjUXsH6RpctkQlwu2VcHfo6S5Hl1lZemk0pHDJZnGFQ2lNidZLDAEy72++",
	"WSKbXjMzZ36J20C/KhFScN5KvYANB+lRueHwsMSkfpyhC/jR4hxS+GG8LWTcVDH+o/rypUDK3xworBUr",
	"2Gw4WPxcSYs5pLcEvE34sdtqZp8w8xT7s7XGDvRBM6QVZtGVRjvsNSQzOQ71j94zi6VFh9pLvQilhWTA",
	"YW6sEh5SkNqfnW6bJbXHBVpCpNA5sRjM3Syx+H62n3yv8XuNCIC36Yea0bQ0Xe/V+ghHH9bOcyuGMmYa",
	"EsM3hXQ4VguFw9G08u340srsQIKw1MugRC1VpSCdJEmScFBSx2e6rwovywJ/n0OajJJJd56u1AyJj/pk",
	"YU6al/PCCH9xvpXrU8R9QzsHBR060aTij8ppqzzM701zfr8JFEA96OpHTcXewmlNWj2rz4DDeX0OHN7U",
	"b4DDRX0BHN7Wb4HDu/rdT3T6ygoldR2ktmhv3R0+zOj6uRIWTybAIStMdgccjF+i3YHbyZfEoueGoHrp",
	"SZdNASy6F/sL7T2VyuEerYtlJKNkNKFaTYlalBJSOBtNRgkhE34Z5Dyede5XGheupHZBnZjmkMJ7i8Lj",
	"b/hw1ToSEYDOX5l8FYdfe9QhUpRlIbMQO/7k4oBEFoNtPsvJNvygZXrDRJ7TZcc5g1KiX4VDT5PJP0L7",
	"NZBNKw6AigDYg3AsC83LIeybi6rwL4YhmvUAhM6uTZZVlgaF9rhKKWFXHaFMMI0PXbM2vFXCeC3zTRyJ",
	"Aj32BXEd3ndqKIUVCj1aB+ntU76KxFOTm0QNaRBjO8dpnOntkHtbId/pyf5kfOzxfD480Du0xOOPi5bY",
	"ViZanLMVm14TwgUOTOUH9M9jwKK3Eu9fi4PkP521oyLyA/oBFofN9TLPpx7VjXkemeR+rRW+HKHfXf5/",
	"7PKXO5Lo5Njz+fE6/t5Mv275f6Iy9xhJ+8Ua9e/V+iLq5E/5Z44mR7iHz2wLf52PTYOB1BFB5GxujdqT",
	"7tGoJRLMRIs7YH2smiWKwi8JweAH6f0Ss7tf457hr0G/SS7+szLpWMy+Oq4fJSopMBbRtfpqYFPKzd8B",
	"AAD//+3ZbbpXEAAA",
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