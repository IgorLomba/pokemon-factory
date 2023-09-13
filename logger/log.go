package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type Log struct {
	Level        string `env:"LOG_LEVEL,default=info"`
	JSON         bool   `env:"LOG_JSON,default=false"`
	APPComponent string `env:"APP_COMPONENT,default=PokemonFactory"`
}

func Config(logSettings *Log) {
	if logSettings == nil {
		logSettings = &Log{}
	}
	if logSettings.APPComponent == "" {
		logSettings.APPComponent = "PokemonFactory"
	}

	level, err := zerolog.ParseLevel(logSettings.Level)
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	var levelNameHook LevelNameHook
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	if logSettings.JSON {
		log.Logger = zerolog.New(os.Stdout).Hook(levelNameHook).With().Stack().Caller().Str("component", logSettings.APPComponent).Timestamp().Logger()
	} else {
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).Hook(levelNameHook).With().Stack().Caller().Str("component", logSettings.APPComponent).Timestamp().Logger()
	}
	log.Info().Str("log_level", level.String()).Msg("Logs configuration : OK")
}

// -1-tracing 0-debug 1-info 2-warning 3-error 4-emergency 5-critical
var logSeverityMap = map[zerolog.Level]string{
	zerolog.TraceLevel: "TRACE",
	zerolog.DebugLevel: "DEBUG",
	zerolog.InfoLevel:  "INFO",
	zerolog.WarnLevel:  "WARNING",
	zerolog.ErrorLevel: "ERROR",
	zerolog.FatalLevel: "EMERGENCY",
	zerolog.PanicLevel: "CRITICAL",
}

type LevelNameHook struct{}

func (h LevelNameHook) Run(e *zerolog.Event, l zerolog.Level, msg string) {
	if l == zerolog.NoLevel {
		e.Str("severity", zerolog.DebugLevel.String())
		return
	}
	e.Str("severity", logSeverityMap[l])
}
