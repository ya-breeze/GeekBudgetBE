// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Geek Budget - OpenAPI 3.0
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.0.1
 * Contact: ilya.korolev@outlook.com
 */

package goserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// CurrenciesAPIController binds http requests to an api service and writes the service results to the http response
type CurrenciesAPIController struct {
	service      CurrenciesAPIServicer
	errorHandler ErrorHandler
}

// CurrenciesAPIOption for how the controller is set up.
type CurrenciesAPIOption func(*CurrenciesAPIController)

// WithCurrenciesAPIErrorHandler inject ErrorHandler into controller
func WithCurrenciesAPIErrorHandler(h ErrorHandler) CurrenciesAPIOption {
	return func(c *CurrenciesAPIController) {
		c.errorHandler = h
	}
}

// NewCurrenciesAPIController creates a default api controller
func NewCurrenciesAPIController(s CurrenciesAPIServicer, opts ...CurrenciesAPIOption) *CurrenciesAPIController {
	controller := &CurrenciesAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the CurrenciesAPIController
func (c *CurrenciesAPIController) Routes() Routes {
	return Routes{
		"GetCurrencies": Route{
			strings.ToUpper("Get"),
			"/v1/currencies",
			c.GetCurrencies,
		},
		"CreateCurrency": Route{
			strings.ToUpper("Post"),
			"/v1/currencies",
			c.CreateCurrency,
		},
		"UpdateCurrency": Route{
			strings.ToUpper("Put"),
			"/v1/currencies/{id}",
			c.UpdateCurrency,
		},
		"DeleteCurrency": Route{
			strings.ToUpper("Delete"),
			"/v1/currencies/{id}",
			c.DeleteCurrency,
		},
	}
}

// GetCurrencies - get all currencies
func (c *CurrenciesAPIController) GetCurrencies(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.GetCurrencies(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	_ = EncodeJSONResponse(result.Body, &result.Code, w)
}

// CreateCurrency - create new currency
func (c *CurrenciesAPIController) CreateCurrency(w http.ResponseWriter, r *http.Request) {
	currencyNoIdParam := CurrencyNoId{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&currencyNoIdParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertCurrencyNoIdRequired(currencyNoIdParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertCurrencyNoIdConstraints(currencyNoIdParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateCurrency(r.Context(), currencyNoIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	_ = EncodeJSONResponse(result.Body, &result.Code, w)
}

// UpdateCurrency - update currency
func (c *CurrenciesAPIController) UpdateCurrency(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idParam := params["id"]
	if idParam == "" {
		c.errorHandler(w, r, &RequiredError{"id"}, nil)
		return
	}
	currencyNoIdParam := CurrencyNoId{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&currencyNoIdParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertCurrencyNoIdRequired(currencyNoIdParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertCurrencyNoIdConstraints(currencyNoIdParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateCurrency(r.Context(), idParam, currencyNoIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	_ = EncodeJSONResponse(result.Body, &result.Code, w)
}

// DeleteCurrency - delete currency
func (c *CurrenciesAPIController) DeleteCurrency(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idParam := params["id"]
	if idParam == "" {
		c.errorHandler(w, r, &RequiredError{"id"}, nil)
		return
	}
	result, err := c.service.DeleteCurrency(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	_ = EncodeJSONResponse(result.Body, &result.Code, w)
}
