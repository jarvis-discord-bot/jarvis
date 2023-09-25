package applies

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pedromsmoreira/jarvis/internal/configuration"
	"github.com/pedromsmoreira/jarvis/internal/discord"
	"github.com/sirupsen/logrus"
)

type Facade struct {
	gateway    *discord.Gateway
	appliesCfg *configuration.Applies
}

func NewFacade(settings *configuration.Settings, gw *discord.Gateway) *Facade {
	return &Facade{
		gateway:    gw,
		appliesCfg: settings.ApiCfg.AppliesCfg,
	}
}

func (s *Facade) Create(apply *Apply) error {
	embed := MapApplyToEmbed(apply)
	if err := s.gateway.PostMessageEmbed(s.appliesCfg.AppliesChannelID, embed); err != nil {
		return fmt.Errorf("error posting apply message: %w", err)
	}

	chID, err := s.gateway.CreateChannel(
		s.appliesCfg.GuildID, fmt.Sprintf("apply-%v", apply.UserID), s.appliesCfg.ApplyParentCategoryID)
	if err != nil {
		return fmt.Errorf("error creating channel: %w", err)
	}

	applierID, err := s.gateway.GetUserID(s.appliesCfg.GuildID, apply.UserID)
	if err != nil {
		return fmt.Errorf("error getting user id: %w", err)
	}

	err = s.gateway.AddUserToChannel(chID, applierID)
	if err != nil {
		return fmt.Errorf("error adding user to channel: %w", err)
	}

	return s.gateway.SendWelcomeMessage(chID)
}

const embedType = discordgo.EmbedTypeRich

func MapApplyToEmbed(apply *Apply) *discordgo.MessageEmbed {
	fields := make([]*discordgo.MessageEmbedField, len(apply.Questions))
	logrus.WithField("questions_number", len(apply.Questions)).Info("questions received")
	for i := 0; i < len(apply.Questions); i++ {
		logrus.
			WithField("question", apply.Questions[i].Question).
			WithField("answer", apply.Questions[i].Answer).
			Infof("question %d", i)
		f := &discordgo.MessageEmbedField{
			Name:   apply.Questions[i].Question,
			Value:  apply.Questions[i].Answer,
			Inline: false,
		}
		logrus.WithField("field", f).Info("created field")
		fields[i] = f
	}

	return &discordgo.MessageEmbed{
		Title:  apply.Title,
		URL:    apply.URL,
		Type:   embedType,
		Fields: fields,
	}
}
