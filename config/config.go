package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

const (
	// DebugMode indicates service mode is debug.
	DebugMode = "debug"
	// TestMode indicates service mode is test.
	TestMode = "test"
	// ReleaseMode indicates service mode is release.
	ReleaseMode = "release"
)

// Config ...
type Config struct {
	ServiceName string
	Version     string
	Environment string

	HTTPScheme string
	HTTPPort   string

	RPCPort string

	MongoHost     string
	MongoPort     int
	MongoUser     string
	MongoPassword string
	MongoDatabase string

	DefaultOffset string
	DefaultLimit  string
}

// Load ...
func Load() Config {
	if err := godotenv.Load("/app/.env"); err != nil {
		if err := godotenv.Load(".env"); err != nil {
			fmt.Println("No .env file found")
		}
		fmt.Println("No .env file found")
	}

	config := Config{}

	config.ServiceName = cast.ToString(getOrReturnDefaultValue("SERVICE_NAME", "chat_service"))
	config.Version = cast.ToString(getOrReturnDefaultValue("VERSION", "1.0"))
	config.Environment = cast.ToString(getOrReturnDefaultValue("ENVIRONMENT", DebugMode))

	config.RPCPort = cast.ToString(getOrReturnDefaultValue("RPC_PORT", ":9001"))

	config.MongoHost = cast.ToString(getOrReturnDefaultValue("MONGO_HOST", "0.0.0.0"))
	config.MongoPort = cast.ToInt(getOrReturnDefaultValue("MONGO_PORT", "17017"))
	config.MongoUser = cast.ToString(getOrReturnDefaultValue("MONGO_USER", "admin"))
	config.MongoPassword = cast.ToString(getOrReturnDefaultValue("MONGO_PASSWORD", "password"))
	config.MongoDatabase = cast.ToString(getOrReturnDefaultValue("MONGO_DATABASE", "chat_service"))

	config.DefaultOffset = cast.ToString(getOrReturnDefaultValue("DEFAULT_OFFSET", "0"))
	config.DefaultLimit = cast.ToString(getOrReturnDefaultValue("DEFAULT_LIMIT", "100"))

	config.HTTPScheme = cast.ToString(getOrReturnDefaultValue("HTTP_SCHEME", "http"))
	config.HTTPPort = cast.ToString(getOrReturnDefaultValue("HTTP_PORT", ":8181"))

	return config
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)

	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
