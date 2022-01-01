package configuration

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type ServicePortConfig struct {
	GRPCPort int `json:"GRPC_PORT"`
	HTTPPort int `json:"HTTP_PORT"`
}

type LoggingConfig struct {
	LogLevel      int    `json:"LOG_LEVEL"`
	LogTimeFormat string `json:"LOG_TIME_FORMAT"`
}

var servicePortConfig *ServicePortConfig
var loggingConfig *LoggingConfig

func PortConfig() ServicePortConfig { return *servicePortConfig }
func LogConfig() LoggingConfig      { return *loggingConfig }

func LoadConfig() {
	file, err := os.Open("./configuration/config.json")

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	byteValue, _ := ioutil.ReadAll(file)

	_ = json.Unmarshal(byteValue, &servicePortConfig)
	_ = json.Unmarshal(byteValue, &loggingConfig)
}
