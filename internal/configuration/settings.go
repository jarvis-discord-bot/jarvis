package configuration

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Settings struct {
	Srv        *Server
	DiscordCfg *Discord
	ApiCfg     *Api
	Sql        *Sql
}

type Server struct {
	Address string
	Port    int
}

type Discord struct {
	Token string `validate:"required"`
}

type Api struct {
	AppliesCfg *Applies
}

type Sql struct {
	Dsn              string `validate:"required"`
	MigrationPath    string `validate:"required"`
	RequestTimeoutMs int
}

type Applies struct {
	Enabled               bool
	GuildID               string
	ApplyParentCategoryID string
	AppliesChannelID      string
}

func NewSettings() (*Settings, error) {
	viper.SetEnvPrefix("JARVIS")
	viper.AutomaticEnv()
	viper.SetDefault("PORT", 8081)
	viper.SetDefault("REQUEST_TIMEOUT_IN_MS", 5000)
	viper.SetDefault("SQL_REQUEST_TIMEOUT_MS", 10000)

	settings := &Settings{
		Srv: &Server{
			Address: viper.GetString("ADDRESS"),
			Port:    viper.GetInt("PORT"),
		},
		DiscordCfg: &Discord{
			Token: strings.TrimSpace(viper.GetString("BOT_TOKEN")),
		},
		ApiCfg: &Api{
			AppliesCfg: &Applies{
				Enabled:               viper.GetBool("APPLY_ENABLED"),
				GuildID:               viper.GetString("APPLY_GUILD_ID"),
				ApplyParentCategoryID: viper.GetString("APPLY_PARENT_CATEGORY_ID"),
				AppliesChannelID:      viper.GetString("APPLY_CHANNEL_ID"),
			},
		},
		Sql: &Sql{
			MigrationPath:    viper.GetString("SQL_MIGRATION_PATH"),
			Dsn:              viper.GetString("SQL_DSN"),
			RequestTimeoutMs: viper.GetInt("SQL_REQUEST_TIMEOUT_MS"),
		},
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(settings); err != nil {
		return nil, err
	}
	return settings, nil
}
