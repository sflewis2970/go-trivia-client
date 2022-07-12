package config

import (
	"os"
	"testing"
)

const (
	HOSTPORTVAL          string = ":3000"
	TRIVIASERVICENAMEVAL string = "trivia-service"
	TRIVIASERVICEPORTVAL string = ":8080"
)

func TestGetConfig(t *testing.T) {
	// Call Get() for the first time
	cfg := Get()

	if cfg == nil {
		t.Error("After initial Get() call cfg should NOT be nil")
	}
}

func setEnvVars() {
	os.Setenv(HOSTNAME, "")
	os.Setenv(HOSTPORT, HOSTPORTVAL)
	os.Setenv(TRIVIASERVICENAME, TRIVIASERVICENAMEVAL)
	os.Setenv(TRIVIASERVICEPORT, TRIVIASERVICEPORTVAL)
}

func TestGetConfigData(t *testing.T) {
	setEnvVars()

	cfgData, _ := Get().GetData(REFRESH_CONFIG_DATA)
	if cfgData == nil {
		t.Error("After call to Get().GetData() cfgData shoud NOT be nil")
		return
	}

	if cfgData.HostPort != HOSTPORTVAL {
		t.Errorf("Invalid host port value..., got: %s, expected: %s", cfgData.HostPort, HOSTPORTVAL)
		return
	}

	if cfgData.TriviaServiceName != TRIVIASERVICENAMEVAL {
		t.Errorf("Invalid trivia service name value..., got: %s, expected: %s", cfgData.TriviaServiceName, TRIVIASERVICENAMEVAL)
		return
	}

	if cfgData.TriviaServicePort != TRIVIASERVICEPORTVAL {
		t.Errorf("Invalid trivia service port value..., got: %s, expected: %s", cfgData.TriviaServicePort, TRIVIASERVICEPORTVAL)
		return
	}
}
