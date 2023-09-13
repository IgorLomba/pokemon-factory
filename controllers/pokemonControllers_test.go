package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pokemon/models"
	"pokemon/util"

	"github.com/rs/zerolog/log"
)

func serveHTTP(u PokemonController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/pokemon/get") {
			u.GetPokemonsPolling()(w, r)
		}
		if strings.Contains(r.URL.Path, "/api/pokemon/generate") {
			u.GeneratePokemons()(w, r)
		}
		if strings.Contains(r.URL.Path, "/api/pokemon/store") {
			u.LatestPokemon()(w, r)
		}
		if strings.Contains(r.URL.Path, "/api/pokemon/status") {
			u.CheckStatus()(w, r)
		}
	}
}

func TestPokemonsControllerHTTP_GeneratePokemons(t *testing.T) {
	util.BindEnvs()
	u := NewPokemonController()
	server := httptest.NewServer(http.HandlerFunc(u.GeneratePokemons()))
	defer server.Close()

	url := fmt.Sprintf("%s%s/generate/%d", server.URL, PokemonAPIPath, 5)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 but got: %d", resp.StatusCode)
	}
	if resp.Body == nil {
		t.Errorf("expected response body but got nil")
	}
	var responseStruct models.RequestStatus
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&responseStruct); err != nil {
		t.Errorf("error decoding response body: %v", err)
		return
	}

	log.Info().Msgf("Pokemon generate test finished with body : %v", responseStruct)
}
func TestPokemonsControllerHTTP_GetPokemonsPolling(t *testing.T) {
	util.BindEnvs()
	u := NewPokemonController()
	server := httptest.NewServer(http.HandlerFunc(serveHTTP(u)))
	defer server.Close()

	PortWithoutURL := strings.Split(server.URL, ":")[2]
	URLWithoutPort := strings.Split(server.URL, ":")[0] + ":" + strings.Split(server.URL, ":")[1]
	util.Envs[util.EnvHostURL] = URLWithoutPort
	util.Envs[util.EnvPort] = PortWithoutURL

	url := fmt.Sprintf("%s%s/generate/%d", server.URL, PokemonAPIPath, 5)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 but got: %d", resp.StatusCode)
	}
	if resp.Body == nil {
		t.Errorf("expected response body but got nil")
	}
	var responseStruct models.RequestStatus
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&responseStruct); err != nil {
		t.Errorf("error decoding response body: %v", err)
		return
	}

	url = fmt.Sprintf("%s%s/get?id=%s", server.URL, PokemonAPIPath, responseStruct.ID)

	poolResp, err := http.Get(url)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer poolResp.Body.Close()
	var respModel interface{}
	decoder = json.NewDecoder(poolResp.Body)
	if err := decoder.Decode(&respModel); err != nil {
		t.Errorf("error decoding response body: %v", err)
		return
	}
	if poolResp.StatusCode != http.StatusAccepted {
		t.Errorf("expected status 200 but got: %d", poolResp.StatusCode)
	}

	log.Info().Msgf("Pokemon polling test finished with body : %v", respModel)
}

func TestPokemonsControllerHTTP_LatestPokemon(t *testing.T) {
	util.BindEnvs()
	u := NewPokemonController()
	server := httptest.NewServer(http.HandlerFunc(serveHTTP(u)))
	defer server.Close()

	PortWithoutURL := strings.Split(server.URL, ":")[2]
	URLWithoutPort := strings.Split(server.URL, ":")[0] + ":" + strings.Split(server.URL, ":")[1]
	util.Envs[util.EnvHostURL] = URLWithoutPort
	util.Envs[util.EnvPort] = PortWithoutURL

	url := fmt.Sprintf("%s%s/generate/%d", server.URL, PokemonAPIPath, 5)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 but got: %d", resp.StatusCode)
	}
	if resp.Body == nil {
		t.Errorf("expected response body but got nil")
	}
	var responseStruct models.RequestStatus
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&responseStruct); err != nil {
		t.Errorf("error decoding response body: %v", err)
		return
	}

	url = fmt.Sprintf("%s%s/get?id=%s", server.URL, PokemonAPIPath, responseStruct.ID)

	pollResp, err := http.Get(url)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer pollResp.Body.Close()
	var polling pokemonControllerHTTP
	decoder = json.NewDecoder(pollResp.Body)
	if err := decoder.Decode(&polling); err != nil {
		t.Errorf("error decoding response body: %v", err)
		return
	}
	if pollResp.StatusCode != http.StatusAccepted {
		t.Errorf("expected status 200 but got: %d", pollResp.StatusCode)
	}

	url = fmt.Sprintf("%s%s/store", server.URL, PokemonAPIPath)
	lastResp, err := http.Get(url)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer lastResp.Body.Close()

	// to test with no pokemons in the list
	if lastResp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404 but got: %d", lastResp.StatusCode)
	}

	pokemonsList.PushBack(polling.Pokemons)

	lastResp, err = http.Get(url)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer lastResp.Body.Close()

	// and to test with pokemons in the list
	if lastResp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 but got: %d", lastResp.StatusCode)
	}

	var last []models.Pokemon
	decoder = json.NewDecoder(lastResp.Body)
	if err := decoder.Decode(&last); err != nil {
		t.Errorf("error decoding response body: %v", err)
		return
	}

	log.Info().Msgf("Pokemon store test finished with body : %v", last)
}
