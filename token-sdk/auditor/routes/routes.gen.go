// Package auditor provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.4 DO NOT EDIT.
package routes

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// Account Information about an account and its balance
type Account struct {
	// Balance balance in base units for each currency
	Balance []Amount `json:"balance"`

	// Id account id as registered at the Certificate Authority
	Id string `json:"id"`
}

// Amount The amount to issue, transfer or redeem.
type Amount struct {
	// Code the code of the token
	Code string `json:"code"`

	// Value value in base units (usually cents)
	Value int64 `json:"value"`
}

// Error defines model for Error.
type Error struct {
	// Message High level error message
	Message string `json:"message"`

	// Payload Details about the error
	Payload string `json:"payload"`
}

// TransactionRecord A transaction
type TransactionRecord struct {
	// Amount The amount to issue, transfer or redeem.
	Amount Amount `json:"amount"`

	// Id transaction id
	Id string `json:"id"`

	// Message user provided message
	Message string `json:"message"`

	// Recipient the recipient of the transaction
	Recipient string `json:"recipient"`

	// Sender the sender of the transaction
	Sender string `json:"sender"`

	// Status Unknown | Pending | Confirmed | Deleted
	Status string `json:"status"`

	// Timestamp timestamp in the format: "2018-03-20T09:12:28Z"
	Timestamp time.Time `json:"timestamp"`
}

// Code The token code to filter on
type Code = string

// Id account id as registered at the Certificate Authority
type Id = string

// AccountSuccess defines model for AccountSuccess.
type AccountSuccess struct {
	Message string `json:"message"`

	// Payload Information about an account and its balance
	Payload Account `json:"payload"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse = Error

// HealthSuccess defines model for HealthSuccess.
type HealthSuccess struct {
	// Message ok
	Message string `json:"message"`
}

// TransactionsSuccess defines model for TransactionsSuccess.
type TransactionsSuccess struct {
	Message string              `json:"message"`
	Payload []TransactionRecord `json:"payload"`
}

// AuditorAccountParams defines parameters for AuditorAccount.
type AuditorAccountParams struct {
	Code *Code `form:"code,omitempty" json:"code,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get an account and their balance of a certain type
	// (GET /auditor/accounts/{id})
	AuditorAccount(ctx echo.Context, id Id, params AuditorAccountParams) error
	// Get all transactions for an account
	// (GET /auditor/accounts/{id}/transactions)
	AuditorTransactions(ctx echo.Context, id Id) error

	// (GET /healthz)
	Healthz(ctx echo.Context) error

	// (GET /readyz)
	Readyz(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// AuditorAccount converts echo context to params.
func (w *ServerInterfaceWrapper) AuditorAccount(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id Id

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params AuditorAccountParams
	// ------------- Optional query parameter "code" -------------

	err = runtime.BindQueryParameter("form", true, false, "code", ctx.QueryParams(), &params.Code)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter code: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AuditorAccount(ctx, id, params)
	return err
}

// AuditorTransactions converts echo context to params.
func (w *ServerInterfaceWrapper) AuditorTransactions(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id Id

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AuditorTransactions(ctx, id)
	return err
}

// Healthz converts echo context to params.
func (w *ServerInterfaceWrapper) Healthz(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Healthz(ctx)
	return err
}

// Readyz converts echo context to params.
func (w *ServerInterfaceWrapper) Readyz(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Readyz(ctx)
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
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/auditor/accounts/:id", wrapper.AuditorAccount)
	router.GET(baseURL+"/auditor/accounts/:id/transactions", wrapper.AuditorTransactions)
	router.GET(baseURL+"/healthz", wrapper.Healthz)
	router.GET(baseURL+"/readyz", wrapper.Readyz)

}

type AccountSuccessJSONResponse struct {
	Message string `json:"message"`

	// Payload Information about an account and its balance
	Payload Account `json:"payload"`
}

type ErrorResponseJSONResponse Error

type HealthSuccessJSONResponse struct {
	// Message ok
	Message string `json:"message"`
}

type TransactionsSuccessJSONResponse struct {
	Message string              `json:"message"`
	Payload []TransactionRecord `json:"payload"`
}

type AuditorAccountRequestObject struct {
	Id     Id `json:"id"`
	Params AuditorAccountParams
}

type AuditorAccountResponseObject interface {
	VisitAuditorAccountResponse(w http.ResponseWriter) error
}

type AuditorAccount200JSONResponse struct{ AccountSuccessJSONResponse }

func (response AuditorAccount200JSONResponse) VisitAuditorAccountResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type AuditorAccountdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response AuditorAccountdefaultJSONResponse) VisitAuditorAccountResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type AuditorTransactionsRequestObject struct {
	Id Id `json:"id"`
}

type AuditorTransactionsResponseObject interface {
	VisitAuditorTransactionsResponse(w http.ResponseWriter) error
}

type AuditorTransactions200JSONResponse struct {
	TransactionsSuccessJSONResponse
}

func (response AuditorTransactions200JSONResponse) VisitAuditorTransactionsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type AuditorTransactionsdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response AuditorTransactionsdefaultJSONResponse) VisitAuditorTransactionsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type HealthzRequestObject struct {
}

type HealthzResponseObject interface {
	VisitHealthzResponse(w http.ResponseWriter) error
}

type Healthz200JSONResponse struct{ HealthSuccessJSONResponse }

func (response Healthz200JSONResponse) VisitHealthzResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type Healthz503JSONResponse struct{ ErrorResponseJSONResponse }

func (response Healthz503JSONResponse) VisitHealthzResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(503)

	return json.NewEncoder(w).Encode(response)
}

type ReadyzRequestObject struct {
}

type ReadyzResponseObject interface {
	VisitReadyzResponse(w http.ResponseWriter) error
}

type Readyz200JSONResponse struct{ HealthSuccessJSONResponse }

func (response Readyz200JSONResponse) VisitReadyzResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type Readyz503JSONResponse struct{ ErrorResponseJSONResponse }

func (response Readyz503JSONResponse) VisitReadyzResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(503)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Get an account and their balance of a certain type
	// (GET /auditor/accounts/{id})
	AuditorAccount(ctx context.Context, request AuditorAccountRequestObject) (AuditorAccountResponseObject, error)
	// Get all transactions for an account
	// (GET /auditor/accounts/{id}/transactions)
	AuditorTransactions(ctx context.Context, request AuditorTransactionsRequestObject) (AuditorTransactionsResponseObject, error)

	// (GET /healthz)
	Healthz(ctx context.Context, request HealthzRequestObject) (HealthzResponseObject, error)

	// (GET /readyz)
	Readyz(ctx context.Context, request ReadyzRequestObject) (ReadyzResponseObject, error)
}

type StrictHandlerFunc = runtime.StrictEchoHandlerFunc
type StrictMiddlewareFunc = runtime.StrictEchoMiddlewareFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// AuditorAccount operation middleware
func (sh *strictHandler) AuditorAccount(ctx echo.Context, id Id, params AuditorAccountParams) error {
	var request AuditorAccountRequestObject

	request.Id = id
	request.Params = params

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.AuditorAccount(ctx.Request().Context(), request.(AuditorAccountRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "AuditorAccount")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(AuditorAccountResponseObject); ok {
		return validResponse.VisitAuditorAccountResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// AuditorTransactions operation middleware
func (sh *strictHandler) AuditorTransactions(ctx echo.Context, id Id) error {
	var request AuditorTransactionsRequestObject

	request.Id = id

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.AuditorTransactions(ctx.Request().Context(), request.(AuditorTransactionsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "AuditorTransactions")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(AuditorTransactionsResponseObject); ok {
		return validResponse.VisitAuditorTransactionsResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// Healthz operation middleware
func (sh *strictHandler) Healthz(ctx echo.Context) error {
	var request HealthzRequestObject

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.Healthz(ctx.Request().Context(), request.(HealthzRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Healthz")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(HealthzResponseObject); ok {
		return validResponse.VisitHealthzResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// Readyz operation middleware
func (sh *strictHandler) Readyz(ctx echo.Context) error {
	var request ReadyzRequestObject

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.Readyz(ctx.Request().Context(), request.(ReadyzRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Readyz")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(ReadyzResponseObject); ok {
		return validResponse.VisitReadyzResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
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