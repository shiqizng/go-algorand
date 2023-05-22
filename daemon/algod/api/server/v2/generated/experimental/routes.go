// Package experimental provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/algorand/oapi-codegen DO NOT EDIT.
package experimental

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	. "github.com/algorand/go-algorand/daemon/algod/api/server/v2/generated/model"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Returns OK if experimental API is enabled.
	// (GET /v2/experimental)
	ExperimentalCheck(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ExperimentalCheck converts echo context to params.
func (w *ServerInterfaceWrapper) ExperimentalCheck(ctx echo.Context) error {
	var err error

	ctx.Set(Api_keyScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ExperimentalCheck(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface, m ...echo.MiddlewareFunc) {
	RegisterHandlersWithBaseURL(router, si, "", m...)
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string, m ...echo.MiddlewareFunc) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/v2/experimental", wrapper.ExperimentalCheck, m...)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
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
	var res = make(map[string]func() ([]byte, error))
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
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
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
