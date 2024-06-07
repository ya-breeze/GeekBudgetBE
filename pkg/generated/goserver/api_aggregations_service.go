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
	"context"
	"net/http"
	"errors"
	"time"
)

// AggregationsAPIService is a service that implements the logic for the AggregationsAPIServicer
// This service should implement the business logic for every endpoint for the AggregationsAPI API.
// Include any external packages or services that will be required by this service.
type AggregationsAPIService struct {
}

// NewAggregationsAPIService creates a default api service
func NewAggregationsAPIService() *AggregationsAPIService {
	return &AggregationsAPIService{}
}

// GetBalances - get balance for filtered transactions
func (s *AggregationsAPIService) GetBalances(ctx context.Context, from time.Time, to time.Time, outputCurrencyID string) (ImplResponse, error) {
	// TODO - update GetBalances with the required logic for this service method.
	// Add api_aggregations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, Aggregation{}) or use other options such as http.Ok ...
	// return Response(200, Aggregation{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("GetBalances method not implemented")
}

// GetExpenses - get expenses for filtered transactions
func (s *AggregationsAPIService) GetExpenses(ctx context.Context, from time.Time, to time.Time, outputCurrencyID string) (ImplResponse, error) {
	// TODO - update GetExpenses with the required logic for this service method.
	// Add api_aggregations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, Aggregation{}) or use other options such as http.Ok ...
	// return Response(200, Aggregation{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("GetExpenses method not implemented")
}

// GetIncomes - get incomes for filtered transactions
func (s *AggregationsAPIService) GetIncomes(ctx context.Context, from time.Time, to time.Time, outputCurrencyID string) (ImplResponse, error) {
	// TODO - update GetIncomes with the required logic for this service method.
	// Add api_aggregations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, Aggregation{}) or use other options such as http.Ok ...
	// return Response(200, Aggregation{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("GetIncomes method not implemented")
}