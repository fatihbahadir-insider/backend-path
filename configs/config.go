package configs

import (
	"backend-path/app/tracing"
	"os"
)

type Config struct {}

func Setup() {
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "backend-path-app"
	}
	serviceVersion := os.Getenv("SERVICE_VERSION")
	if serviceVersion == "" {
		serviceVersion = "1.0.0"
	}
	_ = tracing.Init(serviceName, serviceVersion)

	config := Config{}
	config.GormDatabase()
	config.RedisConfig()
}