package service

import (
	"net/http"

	"github.com/PureML-Inc/PureML/server/models"
)

// HealthCheck godoc
//
//	@Summary		Show the status of server.
//	@Description	Get the status of server.
//	@Tags			Root
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/health [get]
func HealthCheck(request *models.Request) *models.Response {
	return models.NewDataResponse(http.StatusOK, nil, "Server is up and running🚀🎉")
}
