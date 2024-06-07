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
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// NotificationsAPIController binds http requests to an api service and writes the service results to the http response
type NotificationsAPIController struct {
	service NotificationsAPIServicer
	errorHandler ErrorHandler
}

// NotificationsAPIOption for how the controller is set up.
type NotificationsAPIOption func(*NotificationsAPIController)

// WithNotificationsAPIErrorHandler inject ErrorHandler into controller
func WithNotificationsAPIErrorHandler(h ErrorHandler) NotificationsAPIOption {
	return func(c *NotificationsAPIController) {
		c.errorHandler = h
	}
}

// NewNotificationsAPIController creates a default api controller
func NewNotificationsAPIController(s NotificationsAPIServicer, opts ...NotificationsAPIOption) *NotificationsAPIController {
	controller := &NotificationsAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the NotificationsAPIController
func (c *NotificationsAPIController) Routes() Routes {
	return Routes{
		"DeleteNotification": Route{
			strings.ToUpper("Delete"),
			"/v1/notifications/{id}",
			c.DeleteNotification,
		},
		"GetNotifications": Route{
			strings.ToUpper("Get"),
			"/v1/notifications",
			c.GetNotifications,
		},
	}
}

// DeleteNotification - delete notification
func (c *NotificationsAPIController) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idParam := params["id"]
	if idParam == "" {
		c.errorHandler(w, r, &RequiredError{"id"}, nil)
		return
	}
	result, err := c.service.DeleteNotification(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	_ = EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetNotifications - return all notifications
func (c *NotificationsAPIController) GetNotifications(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.GetNotifications(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	_ = EncodeJSONResponse(result.Body, &result.Code, w)
}