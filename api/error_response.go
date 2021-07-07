package api

import (
	"devportal/config"
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	StatusCode   int    `json:"statusCode"`
	ResponseCode string `json:"responseCode"`
	Message      string `json:"message"`
}

func WriteErrorResponse(w http.ResponseWriter, errorResponse ErrorResponse, logMessage string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if logMessage != "" {
		config.AddLogFields(config.Logger).Errorf(logMessage)
	} else {
		config.AddLogFields(config.Logger).Errorf(errorResponse.Message)
	}

	marshalledErrorResponse, _ := json.Marshal(errorResponse)

	w.WriteHeader(errorResponse.StatusCode)
	_, _ = w.Write(marshalledErrorResponse)
}
