package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pedromsmoreira/jarvis/internal/api"
	"github.com/pedromsmoreira/jarvis/internal/chat"
	"github.com/pedromsmoreira/jarvis/internal/configuration"
	"github.com/pedromsmoreira/jarvis/internal/database"
	"github.com/pedromsmoreira/jarvis/internal/discord"
	"github.com/pedromsmoreira/jarvis/internal/logger"
	"github.com/sirupsen/logrus"
)

const TargetMigrationVersion = 1

func main() {
	code := run()
	os.Exit(code)
}

func startBot(dg *discord.Gateway, session *database.SqlSession) error {
	dl := chat.NewReactionHandlerFactory(session)
	dg.AddHandler(dl.BuildMessageCreateHandler())
	return dg.Open()
}

func run() int {
	log := logrus.StandardLogger()
	//TODO: move to env vars
	log.SetLevel(logger.WithLogLevel("info"))
	log.SetFormatter(logger.WithFormatter(logger.FormatterTypeText))
	log.SetOutput(logger.WithOutput("stdout"))

	settings, err := configuration.NewSettings()
	if err != nil {
		logrus.
			WithError(err).
			Error("Error loading settings.")
		return 1
	}

	sqlSession, err := database.Connect(settings.Sql)
	if err != nil {
		logrus.WithError(err).Error("Error initializing SQL Session.")
		return 1
	}
	defer sqlSession.Close()
	logrus.Info("Successfully connected to SQL Database.")

	if !database.DoMigrations(TargetMigrationVersion, settings.Sql, sqlSession.Db) {
		return 1
	}

	gw, err := discord.NewDiscordGateway(settings.DiscordCfg)
	if err != nil {
		logrus.WithError(err).
			Error("Error initializing discord gateway connection.")
		return 1
	}

	err = startBot(gw, sqlSession)
	if err != nil {
		logrus.WithError(err).
			Error("Error opening connection to discord gateway.")
		return 1
	}

	logrus.Info("Discord session opened.")

	defer func(gw *discord.Gateway) {
		err := gw.Close()
		if err != nil {
			logrus.WithError(err).
				Error("Error closing discord gateway connection.")
		}
	}(gw)

	quit := make(chan os.Signal, 1)
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	router := api.NewRouter(settings, gw, sqlSession)
	server := api.NewServer(settings, router)
	addr, err := server.Start(wg)
	if err != nil {
		logrus.
			WithField("address", addr).
			WithError(err).
			Error("Error starting http server.")
		return 1
	}
	logrus.WithField("address", addr).Info("HTTP Server is now listening.")

	defer func() {
		logrus.Info("Shutting down HTTP server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = server.Shutdown(ctx); err != nil {
			logrus.WithError(err).
				Error("HTTP Server shutdown failed.")
		}
	}()

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	logrus.Info("Jarvis BOT is now running. Ctrl-C to exit.")
	<-quit
	logrus.Info("Shutdown requested. Shutting down Jarvis BOT...")
	return 0
}
