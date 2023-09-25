package api

import (
	"github.com/gorilla/mux"
	"github.com/pedromsmoreira/jarvis/internal/api/applies"
	"github.com/pedromsmoreira/jarvis/internal/api/common"
	"github.com/pedromsmoreira/jarvis/internal/api/ping"
	"github.com/pedromsmoreira/jarvis/internal/configuration"
	"github.com/pedromsmoreira/jarvis/internal/database"
	"github.com/pedromsmoreira/jarvis/internal/discord"
)

func NewRouter(settings *configuration.Settings, gateway *discord.Gateway, sqlSession *database.SqlSession) *mux.Router {
	facade := applies.NewFacade(settings, gateway)
	router := mux.NewRouter()

	router.Use(common.LoggedHandler)
	router.HandleFunc("/ping", ping.BuildHandler(sqlSession)).Methods("GET")
	if settings.ApiCfg.AppliesCfg.Enabled {
		router.HandleFunc("/applies", applies.BuildHandler(facade)).Methods("POST")
	}

	return router
}
