package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pedromsmoreira/jarvis/internal/configuration"
	"github.com/sirupsen/logrus"
)

// Gateway struct which allows the use of Discord API
type Gateway struct {
	session *discordgo.Session
	cfg     *configuration.Discord
}

var (
	allowApplyPermissionValue int64 = 3072
	denyApplyPermissionValue  int64 = 0
	defaultWelcomeMessage           = "Boas, obrigado pelo apply. Vamos utilizar este canal para te fazer algumas questões, se necessárias."
)

// NewDiscordGateway method should create new DiscordGateway
func NewDiscordGateway(cfg *configuration.Discord) (*Gateway, error) {
	configureLogger()
	session, err := discordgo.New(fmt.Sprintf("Bot %s", cfg.Token))
	if err != nil {
		return nil, err
	}

	return &Gateway{
		session: session,
		cfg:     cfg,
	}, err
}

func (dg *Gateway) PostMessageEmbed(channelID string, msg *discordgo.MessageEmbed) error {
	_, err := dg.session.ChannelMessageSendEmbed(channelID, msg)
	if err != nil {
		return fmt.Errorf("discordgo embed send failed: %w", err)
	}
	return nil
}

func (dg *Gateway) CreateChannel(guildID string, channelName string, parentId string) (channelId string, err error) {
	chData := discordgo.GuildChannelCreateData{
		Name:     channelName,
		ParentID: parentId,
	}

	ch, err := dg.session.GuildChannelCreateComplex(guildID, chData)
	if err != nil {
		return "", err
	}
	return ch.ID, nil
}

func (dg *Gateway) GetUserID(guildID string, userID string) (string, error) {
	members, err := dg.session.GuildMembers(guildID, "", 1000)
	if err != nil {
		return "", fmt.Errorf("error getting guild members: %w", err)
	}

	for i := 0; i < len(members); i++ {
		if userID == members[i].User.Username {
			return members[i].User.ID, nil
		}
	}

	return "", fmt.Errorf("user %v does not exist in Guild Discord %v", userID, guildID)
}

func (dg *Gateway) AddUserToChannel(chID string, userID string) error {
	return dg.session.ChannelPermissionSet(chID, userID, discordgo.PermissionOverwriteTypeMember, allowApplyPermissionValue, denyApplyPermissionValue)
}

func (dg *Gateway) SendWelcomeMessage(chID string) error {
	_, err := dg.session.ChannelMessageSend(chID, defaultWelcomeMessage)
	return err
}

func (dg *Gateway) AddHandler(handler interface{}) {
	dg.session.AddHandler(handler)
}

func (dg *Gateway) Open() error {
	return dg.session.Open()
}

func (dg *Gateway) Close() error {
	return dg.session.Close()
}

func configureLogger() {
	discordgo.Logger = func(msgL, caller int, format string, a ...interface{}) {
		switch msgL {
		case discordgo.LogError:
			logrus.Errorf(format, a...)
		case discordgo.LogInformational:
			logrus.Infof(format, a...)
		case discordgo.LogWarning:
			logrus.Warnf(format, a...)
		case discordgo.LogDebug:
		default:
			logrus.Warnf("unrecognized log level %v. LogEntry: "+format, append([]interface{}{msgL}, a...)...)
		}
	}
}
