package api

import (
	"NemWebGoApi/api/controllers"
	"NemWebGoApi/internal/config"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var server = controllers.Server{}
var cfg *config.Config

func Init(testing bool) controllers.Server {
	if err := godotenv.Load(); err != nil {
		log.Warnln("Error reading .env file: ", err)
	}

	cfg = config.New()

	err := server.Init(
		cfg.SQLFilePath(),
	)

	if err != nil {
		log.Fatalf("server.Initialize: problem initing server: %v", err)
	}
	return server
}

func Run(testing bool) {
	server.Run(":" + cfg.Port())
}
