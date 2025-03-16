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
	"errors"
	"net/http"
)

// MatchersAPIService is an interface that defines the logic for the MatchersAPIServicer
type MatchersAPIService interface {
	// GetMatchers - get all matchers
	GetMatchers(ctx context.Context) (ImplResponse, error)
	// CreateMatcher - create new matcher
	CreateMatcher(ctx context.Context, matcherNoId MatcherNoId) (ImplResponse, error)
	// UpdateMatcher - update matcher
	UpdateMatcher(ctx context.Context, id string, matcherNoId MatcherNoId) (ImplResponse, error)
	// DeleteMatcher - delete matcher
	DeleteMatcher(ctx context.Context, id string) (ImplResponse, error)
	// CheckMatcher - check if passed matcher matches given transaction
	CheckMatcher(ctx context.Context, checkMatcherRequest CheckMatcherRequest) (ImplResponse, error)
}

// MatchersAPIService is a service that implements the logic for the MatchersAPIServicer
// This service should implement the business logic for every endpoint for the MatchersAPI API.
// Include any external packages or services that will be required by this service.
type MatchersAPIServiceImpl struct {
}

// NewMatchersAPIService creates a default api service
func NewMatchersAPIService() MatchersAPIService {
	return &MatchersAPIServiceImpl{}
}

// GetMatchers - get all matchers
func (s *MatchersAPIServiceImpl) GetMatchers(ctx context.Context) (ImplResponse, error) {
	// TODO - update GetMatchers with the required logic for this service method.
	// Add api_matchers_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, []Matcher{}) or use other options such as http.Ok ...
	// return Response(200, []Matcher{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("GetMatchers method not implemented")
}

// CreateMatcher - create new matcher
func (s *MatchersAPIServiceImpl) CreateMatcher(ctx context.Context, matcherNoId MatcherNoId) (ImplResponse, error) {
	// TODO - update CreateMatcher with the required logic for this service method.
	// Add api_matchers_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, Matcher{}) or use other options such as http.Ok ...
	// return Response(200, Matcher{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("CreateMatcher method not implemented")
}

// UpdateMatcher - update matcher
func (s *MatchersAPIServiceImpl) UpdateMatcher(ctx context.Context, id string, matcherNoId MatcherNoId) (ImplResponse, error) {
	// TODO - update UpdateMatcher with the required logic for this service method.
	// Add api_matchers_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, Matcher{}) or use other options such as http.Ok ...
	// return Response(200, Matcher{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("UpdateMatcher method not implemented")
}

// DeleteMatcher - delete matcher
func (s *MatchersAPIServiceImpl) DeleteMatcher(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update DeleteMatcher with the required logic for this service method.
	// Add api_matchers_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	// return Response(200, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("DeleteMatcher method not implemented")
}

// CheckMatcher - check if passed matcher matches given transaction
func (s *MatchersAPIServiceImpl) CheckMatcher(ctx context.Context, checkMatcherRequest CheckMatcherRequest) (ImplResponse, error) {
	// TODO - update CheckMatcher with the required logic for this service method.
	// Add api_matchers_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, CheckMatcher200Response{}) or use other options such as http.Ok ...
	// return Response(200, CheckMatcher200Response{}), nil

	// TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	// return Response(400, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("CheckMatcher method not implemented")
}
