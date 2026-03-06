package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func ReadConfiguration() {

	// load env vars from .env file (optional, ignored if not present).
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			logrus.Warnf("Failed to load the env vars: %v", err)
		}
	}

	// configure logrus
	configureLogrus()
}

// GetEnvString retrieves the value of the environment variable named by the key
// It trims any surrounding whitespace and quotes from the value before returning it.
func GetEnvString(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	value = strings.Trim(value, "\"'")

	if len(value) == 0 {
		logrus.Errorf("Environment variable %s is not set or empty", key)
	}

	return value
}

// configureLogrus sets the log level and formatter for logrus based on the LOGLEVEL environment variable.
func configureLogrus() {
	askedLoglevel := GetEnvString("LOGLEVEL")

	switch askedLoglevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.Warnf("Invalid LOGLEVEL %q, defaulting to debug", askedLoglevel)
		logrus.SetLevel(logrus.DebugLevel)
	}

	// set logrus to output plain text without timestamps or colors,
	// as the logs will be collected by OpenTelemetry and formatted there.
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    true,
	})

}
