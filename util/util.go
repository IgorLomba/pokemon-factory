package util

import (
	"fmt"
	"os"
	"strings"

	"pokemon/logger"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func GenerateUUID4() string {
	uuID4 := uuid.New()
	return fmt.Sprintf("%s", uuID4)
}

var Envs = make(map[string]string)

func BindEnvs() {
	err := os.Setenv(EnvPort, getEnvOrDefault(EnvPort, defaultAPIPort))
	if err != nil {
		log.Err(err).Msg("error setting port env")
		return
	}
	err = os.Setenv(EnvHostURL, getEnvOrDefault(EnvHostURL, localHost))
	if err != nil {
		log.Err(err).Msg("error setting host env")
		return
	}
	configureLogger()
	return
}

func getEnvOrDefault(envVar, defaultValue string) string {
	value := os.Getenv(envVar)
	if value == "" {
		value = defaultValue
	}
	Envs[envVar] = value
	return value
}

func configureLogger() {
	logger.Config(&logger.Log{
		Level:        getEnvOrDefault(envLogLevel, zerolog.DebugLevel.String()),
		JSON:         strings.ToLower(os.Getenv(envLogJson)) == "true",
		APPComponent: getEnvOrDefault(envAPPComponent, defaultAPPComponent),
	})
}
