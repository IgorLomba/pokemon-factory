package logger

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		_ = os.Setenv("APP_COMPONENT", "")
	}()

	err := os.Setenv("APP_COMPONENT", "test_component")
	if err != nil {
		t.Error(err)
		return
	}

	Config(&Log{
		Level:        "debug",
		JSON:         true,
		APPComponent: "test_component",
	})

	err = w.Close()
	if err != nil {
		return
	}
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	var logOutput map[string]string
	err = json.Unmarshal(out, &logOutput)
	if err != nil {
		t.Errorf("error parsing log output: %s", err.Error())
		return
	}
	if logOutput["level"] != "info" {
		t.Errorf("we were expecting the level prop to had an 'info' value but we got : %s", logOutput["level"])
		return
	}
	if logOutput["component"] != "test_component" {
		t.Errorf("we were expecting the component prop to had an 'test_component' value but we got : %s", logOutput["component"])
		return
	}
	if logOutput["log_level"] != "debug" {
		t.Errorf("we were expecting the log_level prop to had an 'debug' value but we got : %s", logOutput["log_level"])
		return
	}
	if logOutput["severity"] != "INFO" {
		t.Errorf("we were expecting the severity prop to had an 'INFO' value but we got : %s", logOutput["severity"])
		return
	}
	if logOutput["message"] != "Logs configuration : OK" {
		t.Errorf("we were expecting the message prop to had an 'Logs configuration : OK' value but we got : %s", logOutput["message"])
		return
	}
}
