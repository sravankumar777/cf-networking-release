package handlers

import (
	"net/http"
)

type DestinationsIndex struct {
	ErrorResponse errorResponse
}

func NewDestinationsIndex (errorResponse errorResponse) *DestinationsIndex{
	return &DestinationsIndex{
		ErrorResponse: errorResponse,
	}
}

func (d *DestinationsIndex) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}