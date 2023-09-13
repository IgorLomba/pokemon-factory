package services

import (
	"net/http"

	"pokemon/models"

	"github.com/rs/zerolog/log"
)

const (
	errIDNotFound          = "ID not found"
	errCreatingRequest     = "error creating request"
	errNoPokemonsAvailable = "no pokemons available"
	errEncodingJSON        = "error encoding JSON"
	errWritingResponse     = "error writing response"
)

func HandleInternalServerError(w http.ResponseWriter) {
	genericError(w, errCreatingRequest, http.StatusInternalServerError, true)
}

func HandleNotFoundError(w http.ResponseWriter, printLog bool) {
	genericError(w, errIDNotFound, http.StatusNotFound, printLog)
}

func HandleNoPokemonsAvailable(w http.ResponseWriter) {
	log.Warn().Msg(errNoPokemonsAvailable)
	genericError(w, errNoPokemonsAvailable, http.StatusNotFound, false)
}

func HandleEncodingJSONError(w http.ResponseWriter) {
	genericError(w, errEncodingJSON, http.StatusInternalServerError, true)
}

func HandleErrorWritingResponse(w http.ResponseWriter) {
	genericError(w, errWritingResponse, http.StatusInternalServerError, true)
}

func genericError(w http.ResponseWriter, message string, status int, printLog bool) {
	if printLog {
		log.Warn().Msgf("error: %s - status: %d", message, status)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(models.ErrorResponse{
		Status:  status,
		Message: message,
	}.ToJSON())
	if err != nil {
		log.Error().Err(err).Msg("error writing response")
	}
}
