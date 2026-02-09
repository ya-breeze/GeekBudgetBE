package api

import (
	"net/http"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/version"
)

// StatusResponse represents the response from the status endpoint.
type StatusResponse struct {
	BuildTime string `json:"buildTime"`
	Commit    string `json:"commit"`
	StartTime string `json:"startTime"`
}

// StatusAPIController handles status-related requests.
type StatusAPIController struct{}

// NewStatusAPIController creates a new StatusAPIController.
func NewStatusAPIController() *StatusAPIController {
	return &StatusAPIController{}
}

// Routes returns the routes for the status API.
func (c *StatusAPIController) Routes() goserver.Routes {
	return goserver.Routes{
		"GetStatus": goserver.Route{
			Method:      http.MethodGet,
			Pattern:     "/v1/status",
			HandlerFunc: c.GetStatus,
		},
	}
}

// GetStatus handles GET /v1/status.
func (c *StatusAPIController) GetStatus(w http.ResponseWriter, _ *http.Request) {
	response := StatusResponse{
		BuildTime: version.BuildTime,
		Commit:    version.Commit,
		StartTime: version.StartTime.Format(time.RFC3339),
	}

	_ = goserver.EncodeJSONResponse(response, nil, w)
}
