package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// Config is a struct for keeping all configuration variables
type Config struct {
	sqlitePath   string
	influxURL    string
	influxToken  string
	influxOrg    string
	influxBucket string
	apiPort      string
	logLevel     string
	testing      bool
}

// New Loads a new config
func New() *Config {
	if _, err := os.Stat("/.dockerenv"); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		err = godotenv.Load()
		if err != nil {
			fmt.Printf("Error getting env, %v", err)
			os.Exit(1)
		} else {
			fmt.Println(".env File Read")
			log.Infoln(".env File Read")
		}
	}

	conf := &Config{}

	conf.sqlitePath = parseEnvString("SQLITE_PATH", "/data/database.sqlite")
	conf.influxURL = parseEnvString("INFLUX_URL", "http://localhost:8086")
	conf.influxToken = parseEnvString("INFLUX_TOKEN", "aaaaaaa")
	conf.influxOrg = parseEnvString("INFLUX_ORG", "nema")
	conf.influxBucket = parseEnvString("INFLUX_BUCKET", "nema_bucket")
	conf.apiPort = parseEnvString("API_PORT", "3005")
	conf.logLevel = parseEnvString("LOG_LEVEL", "info")
	conf.testing, _ = strconv.ParseBool(parseEnvString("TESTING", "False"))

	if conf.testing {
		log.Warnln("TESTING")
	}

	setupLogger(conf.logLevel)

	return conf
}

func parseEnvString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func setupLogger(logLevel string) {
	log.SetOutput(os.Stdout)

	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	switch logLevel {
	case "trace":
		log.Warnln("Log Level: Trace")
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.Warnln("Log Level: Debug")
		log.SetLevel(log.DebugLevel)
	case "info":
		log.Warnln("Log Level: Info")
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.Warnln("Log Level: Warn")
		log.SetLevel(log.WarnLevel)
	case "fatal":
		log.Warnln("Log Level: Fatal")
		log.SetLevel(log.FatalLevel)
	}
}

// Testing returns a boolean for testing status
func (c *Config) Testing() bool {
	return c.testing
}

// SQLFilePath returns the file path for the sqlite database
func (c *Config) SQLFilePath() string {
	return c.sqlitePath
}

// Port returns the port for api access
func (c *Config) Port() string {
	return c.apiPort
}

func (c *Config) InfluxHost() string {
	return c.influxURL
}

func (c *Config) InfluxToken() string {
	return c.influxToken
}

func (c *Config) InfluxBucket() string {
	return c.influxBucket
}

func (c *Config) InfluxOrg() string {
	return c.influxOrg
}
