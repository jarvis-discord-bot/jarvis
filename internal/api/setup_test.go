package api

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/pedromsmoreira/jarvis/internal/configuration"
	"github.com/pedromsmoreira/jarvis/internal/discord"

	"github.com/pedromsmoreira/jarvis/internal/logger"

	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	log := logrus.StandardLogger()
	log.SetLevel(logger.WithLogLevel("info"))
	log.SetFormatter(logger.WithFormatter("json"))
	log.SetOutput(logger.WithOutput("stdout"))

	settings, err := configuration.NewSettings()
	if err != nil {
		log.Fatalf("error creating settings: %v", err)
	}

	gw, err := discord.NewDiscordGateway(settings.DiscordCfg)
	if err != nil {
		log.Fatalf("error creating discord gateway: %v", err)
	}
	router := NewRouter(settings, gw)

	server := NewServer(settings, router)
	go func() {
		err = server.Start()
		if err != nil && err == http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	code := m.Run()

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	os.Exit(code)
}
