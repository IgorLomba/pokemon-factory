//	 Pokemon API:
//	  version: 1.0
//	  title: Pokemon API
//	 Schemes: http, https
//	 Host:
//	 BasePath: /
//		Consumes:
//		- application/json
//	 Produces:
//	 - application/json
//	 SecurityDefinitions:
//	  Bearer:
//	   type: apiKey
//	   name: Authorization
//	   in: header
//	 swagger:meta
package main

import (
	"fmt"
	"net/http"

	"pokemon/controllers"
	"pokemon/routes"
	"pokemon/util"

	"github.com/rs/zerolog/log"
)

func main() {
	if err := runServer(); err != nil {
		log.Fatal().Err(err).Msg("error starting server")
	}
}

func runServer() error {
	util.BindEnvs()

	pokemonController := controllers.NewPokemonController()
	routes.RegisterRoutes(pokemonController)

	log.Info().Msgf("listening server at:  %s:%s", util.Envs[util.EnvHostURL], util.Envs[util.EnvPort])
	err := http.ListenAndServe(fmt.Sprintf(":%s", util.Envs[util.EnvPort]), nil)
	if err != nil {
		return err
	}

	log.Info().Msg("stopping pokemon server")
	return nil
}
