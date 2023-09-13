package controllers

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"pokemon/models"
	"pokemon/services"
	"pokemon/util"
)

const (
	StatusInProgress = "in progress"
	statusCompleted  = "completed"
	pollingDelay     = 10
	PokemonAPIPath   = "/api/pokemon"
	ID               = "id"
	applicationsJSON = "application/json"
	contentType      = "Content-Type"
)

var (
	names, adj, capabilities []string
	pokemonsList             = list.New()
	requestQueue             = list.New()
	requestMutex             sync.Mutex
)

func init() {
	names, adj, capabilities = services.GetNameAndAdjFromFiles("pokeNames.txt", "adj.txt", "capabilities.txt")
	go processRequests()
}

type requestInfo struct {
	u         *pokemonControllerHTTP
	r         *http.Request
	requestID string
}

func NewPokemonController() PokemonController {
	return &pokemonControllerHTTP{
		Pokemons:      []models.Pokemon{},
		RequestStatus: &models.RequestStatus{},
	}
}

// CheckStatus check the status of the request
/**
 * swagger:route GET /api/pokemon/status PokemonQueryParameters
 *
 *  Check pokemonsList generation status by ID
 *
 *	Produces:
 *	- application/json
 *
 * 	Responses:
 * 		200: RequestStatus
 *      400: ErrorResponse
 *      500: ErrorResponse
 *
 */
func (u *pokemonControllerHTTP) CheckStatus() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		requestID := values.Get(ID)
		log.Trace().Msgf("checking pokemon status - requestID: %s", requestID)

		if _, ok := pokemonHashMap[requestID]; !ok {
			services.HandleNotFoundError(w, false)
			return
		}
		if pokemonHashMap[requestID].RequestStatus.Completed {
			jsonBody, _ := json.Marshal(pokemonHashMap[requestID].RequestStatus)
			w.Header().Set(contentType, applicationsJSON)
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(jsonBody)
			log.Debug().Bytes("CHECK_STATUS_RESPONSE", jsonBody).Msg("response")
			if err != nil {
				services.HandleErrorWritingResponse(w)
				return
			}
			return
		}
		if pokemonHashMap[requestID].RequestStatus.Status != "" {
			w.Header().Set(contentType, applicationsJSON)
			w.WriteHeader(http.StatusAccepted)
			_, err := w.Write(models.ErrorResponseWithRequestStatus{
				Status:        http.StatusAccepted,
				RequestStatus: *pokemonHashMap[requestID].RequestStatus,
			}.ToJSON())
			if err != nil {
				services.HandleErrorWritingResponse(w)
				return
			}
			return
		}
		return
	}
}

// GetPokemonsPolling would be possible to perform using only the hashmap
/**
 * swagger:route GET /api/pokemon/get PokemonQueryParameters
 *
 *  Get pokemonsList by ID
 *
 *	Produces:
 *	- application/json
 *
 * 	Responses:
 * 		200: RequestStatus
 *      400: ErrorResponse
 *      500: ErrorResponse
 *
 */
func (u *pokemonControllerHTTP) GetPokemonsPolling() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		requestID := values.Get(ID)
		log.Info().Msg("getting pokemonsList by id - requestID: " + requestID)

		statusURL := fmt.Sprintf("%s:%s%s/status?id=%s", util.Envs[util.EnvHostURL], util.Envs[util.EnvPort], PokemonAPIPath, requestID)
		req, err := http.NewRequest(http.MethodGet, statusURL, nil)
		if err != nil {
			services.HandleInternalServerError(w)
			return
		}
		client := &http.Client{
			Timeout: 60 * time.Second,
		}
		resp := &http.Response{
			StatusCode: http.StatusAccepted,
		}

		for resp.StatusCode != http.StatusOK {
			resp, err = client.Do(req)
			if err != nil {
				services.HandleInternalServerError(w)
				return
			}
			if resp.StatusCode == http.StatusNotFound {
				services.HandleNotFoundError(w, true)
				return
			}
			log.Trace().Int("POLLING_STATUS", resp.StatusCode).Msg("status return from polling")
			time.Sleep(pollingDelay * time.Millisecond)
		}
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				log.Error().Err(err).Msg("error closing response body while polling")
				return
			}
		}(resp.Body)

		jsonPokemon, err := json.Marshal(pokemonHashMap[requestID])
		if err != nil {
			return
		}
		log.Debug().Bytes("POLLING_RESPONSE", jsonPokemon).Msg("json response")
		w.Header().Set(contentType, applicationsJSON)
		w.WriteHeader(http.StatusAccepted)
		_, err = w.Write(jsonPokemon)
		if err != nil {
			services.HandleErrorWritingResponse(w)
			return
		}
		return
	}
}

// GeneratePokemons get the last pokemonsList generated
/**
 * swagger:route GET /api/pokemon/generate PokemonGenerateQueryParameters
 *
 *  Generate pokemonsList
 *
 *	Produces:
 *	- application/json
 *
 * 	Responses:
 * 		200: RequestStatus
 *      500: ErrorResponse
 *
 */
func (u *pokemonControllerHTTP) GeneratePokemons() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := util.GenerateUUID4()
		log.Info().Msg("generating pokemonsList - requestID: " + requestID)
		w.Header().Set(contentType, applicationsJSON)
		u = &pokemonControllerHTTP{
			RequestStatus: &models.RequestStatus{
				ID:        requestID,
				Status:    StatusInProgress,
				Completed: false,
			},
		}
		pokemonHashMap[requestID] = *u
		log.Info().Msgf("new request - requestID: %s", requestID)

		requestMutex.Lock()
		requestQueue.PushBack(requestInfo{u: u, r: r, requestID: requestID})
		requestMutex.Unlock()

		jsonBody, _ := json.Marshal(u.RequestStatus)
		_, err := w.Write(jsonBody)
		if err != nil {
			services.HandleErrorWritingResponse(w)
			return
		}
	}
}

// LatestPokemon get the last pokemonsList generated
/**
 * swagger:route GET /api/pokemon/store Store
 *
 *  Get the last pokemonsList generated
 *
 *	Produces:
 *	- application/json
 *
 * 	Responses:
 * 		200: Pokemon
 *      500: ErrorResponse
 *
 */
func (u *pokemonControllerHTTP) LatestPokemon() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msg("getting latest pokemon")
		if pokemonsList.Len() == 0 {
			services.HandleNoPokemonsAvailable(w)
			return
		}

		lastPokemon := pokemonsList.Back().Value.([]models.Pokemon)
		pokemonsList.Remove(pokemonsList.Back())
		jsonBody, err := json.Marshal(lastPokemon)
		if err != nil {
			services.HandleEncodingJSONError(w)
			return
		}

		w.Header().Set(contentType, applicationsJSON)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonBody)
		if err != nil {
			services.HandleErrorWritingResponse(w)
			return
		}
	}
}
func processRequests() {
	for {
		requestMutex.Lock()
		if requestQueue.Len() == 0 {
			requestMutex.Unlock()
			continue
		}
		elem := requestQueue.Front()
		requestQueue.Remove(elem)
		requestMutex.Unlock()

		requestInfo := elem.Value.(requestInfo)
		u := requestInfo.u
		r := requestInfo.r
		requestID := requestInfo.requestID

		processPokemonRequest(u, r, requestID)
	}
}

func processPokemonRequest(u *pokemonControllerHTTP, r *http.Request, requestID string) {
	log.Info().Msgf("processing requestID: %s", requestID)
	sleepTime := time.Duration(rand.Intn(1000)) * time.Millisecond
	values := r.URL.Query()
	amount, _ := strconv.Atoi(values.Get("amount"))

	var items []models.Pokemon
	for j := 1; j < amount; j++ {
		name := adj[rand.Intn(len(adj))] + "-" + names[rand.Intn(len(names))]

		item := models.Pokemon{
			Name: name,
		}
		time.Sleep(sleepTime)

		item.Capabilities = services.AppendCapabilities(item.Capabilities, capabilities)
		items = append(items, item)
	}
	u.Pokemons = items
	u.RequestStatus.Completed = true
	u.RequestStatus.Status = statusCompleted
	log.Info().Msgf("pokemon ready - requestID: %s", requestID)

	if len(u.Pokemons) > 0 {
		pokemonsList.PushBack(u.Pokemons)
	}
	pokemonHashMap[requestID] = *u
}
