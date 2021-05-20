package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	Version               string   `envconfig:"VERSION" required:"false" default:"v0.0.0"`
	Name                  string   `envconfig:"NAME" required:"false" default:"api"`
	Port                  string   `envconfig:"PORT" required:"false" default:"8000"`
	HealthPort            string   `envconfig:"HEALTH_PORT" required:"false" default:"8080"`
	MetricsPort           string   `envconfig:"METRICS_PORT" required:"false" default:"9400"`
	RestPrefix            string   `envconfig:"REST_PREFIX" required:"false" default:"/rest"`
	WebsocketPrefix       string   `envconfig:"WEBSOCKET_PREFIX" required:"false" default:"/ws"`
	HealthPrefix          string   `envconfig:"HEALTH_PREFIX" required:"false" default:"/healthcheck"`
	MetricsPrefix         string   `envconfig:"METRICS_PREFIX" required:"false" default:"/metrics"`
	HealthPollingInterval int      `envconfig:"HEALTH_POLLING_INTERVAL" required:"false" default:"10"`
	LogLevel              string   `envconfig:"LOG_LEVEL" required:"false" default:"INFO"`
	LogToFile             bool     `envconfig:"LOG_TO_FILE" required:"false" default:"false"`
	NetworkName           string   `envconfig:"NETWORK_NAME" required:"false" default:"mainnet"`
	KafkaBrokerURL        string   `envconfig:"KAFKA_BROKER_URL" required:"false" default:""`
	SchemaRegistryURL     string   `envconfig:"SCHEMA_REGISTRY_URL" required:"false" default:""`
	KafkaGroupID          string   `envconfig:"KAFKA_GROUP_ID" required:"false" default:"websocket-group"`
	TopicNames            string   `envconfig:"TOPIC_NAMES" required:"false" default:"blocks"`
	SchemaNames           string   `envconfig:"SCHEMA_NAMES" required:"false" default:"block:block_raw"` // each pair is <topic name>:<topic file in schemas dir>
	ConsumerTopics        []string `envconfig:"CONSUMER_TOPICS" required:"false" default:"[blocks]"`
	ProducerTopics        []string `envconfig:"PRODUCER_TOPICS" required:"false" default:"[blocks]"`
}

var Vars Environment

// Run once on main.go
func GetEnvironment() {
	//Get environment variable file
	env_file := os.Getenv("ENV_FILE")
	if env_file != "" {
		_ = godotenv.Load(env_file)
	} else {
		_ = godotenv.Load()
	}

	err := envconfig.Process("", &Vars)
	if err != nil {
		log.Fatalf("ERROR: envconfig - %s\n", err.Error())
	}

	vars, _ := json.Marshal(Vars)
	log.Info("Config Vars: " + string(vars))
}