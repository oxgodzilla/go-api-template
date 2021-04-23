package config

import (
	"os"

	"testing"
)

func TestEnvironment(t *testing.T) {

	// Set env
	env_map := map[string]string{
		"PORT":             "port",
		"HEALTH_PORT":      "health_port",
		"REST_PREFIX":      "rest_prefix",
		"WEBSOCKET_PREFIX": "websocket_prefix",
		"METRICS_PORT":     "metrics_port",
		"NETWORK_NAME":     "network_name",
	}

	for k, v := range env_map {
		os.Setenv(k, v)
	}

	// Load env
	GetEnvironment()

	// Check env
	if Vars.Port != env_map["PORT"] {
		t.Errorf("Invalid value for env variable: PORT")
	}
	if Vars.HealthPort != env_map["HEALTH_PORT"] {
		t.Errorf("Invalid value for env variable: HEALTH_PORT")
	}
	if Vars.RestPrefix != env_map["REST_PREFIX"] {
		t.Errorf("Invalid value for env variable: REST_PREFIX")
	}
	if Vars.WebsocketPrefix != env_map["WEBSOCKET_PREFIX"] {
		t.Errorf("Invalid value for env variable: WEBSOCKET_PREFIX")
	}
	if Vars.MetricsPort != env_map["METRICS_PORT"] {
		t.Errorf("Invalid value for env variable: METRICS_PORT")
	}
	if Vars.NetworkName != env_map["NETWORK_NAME"] {
		t.Errorf("Invalid value for env variable: NETWORK_NAME")
	}
}
