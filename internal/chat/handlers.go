package chat

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pedromsmoreira/jarvis/internal/database"
	"github.com/sirupsen/logrus"
)

const CommandKeySizeLimit = 50

type ReactionHandlerFactory struct {
	sqlSession *database.SqlSession
}

func NewReactionHandlerFactory(sqlSession *database.SqlSession) *ReactionHandlerFactory {
	return &ReactionHandlerFactory{sqlSession: sqlSession}
}

func (dl *ReactionHandlerFactory) BuildMessageCreateHandler() func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		sqlTimeout := time.Duration(dl.sqlSession.SqlSettings.RequestTimeoutMs) * time.Millisecond

		var re = regexp.MustCompile(`^!(?i:addreaction)[ ]+(\S+)[ ]+(.*)$`)
		matches := re.FindStringSubmatch(m.Content)
		if len(matches) == 3 {
			re = regexp.MustCompile(`^[A-Za-z0-9]+$`)
			if !re.MatchString(matches[1]) {
				_, err := s.ChannelMessageSend(m.ChannelID, "Reaction key can only contain letters and numbers.")
				if err != nil {
					logrus.WithError(err).
						WithField("key", matches[1]).
						WithField("value", matches[2]).
						WithField("channelId", m.ChannelID).
						Error("Error sending reaction validation error message to discord channel.")
				}
				logrus.
					WithField("key", matches[1]).
					WithField("value", matches[2]).
					Info("Could not add reaction because key validation failed.")
				return
			}
			if len(matches[1]) > CommandKeySizeLimit {
				_, err := s.ChannelMessageSend(m.ChannelID, "Reaction key can only contain up to "+strconv.Itoa(CommandKeySizeLimit)+" characters.")
				if err != nil {
					logrus.WithError(err).
						WithField("key", matches[1]).
						WithField("value", matches[2]).
						WithField("channelId", m.ChannelID).
						Error("Error sending reaction size validation error message to discord channel.")
				}
				logrus.
					WithField("key", matches[1]).
					WithField("value", matches[2]).
					Info("Could not add reaction because key size validation failed.")
				return
			}
			logrus.
				WithField("key", matches[1]).
				WithField("value", matches[2]).
				Debug("Adding reaction...")
			ctx, cancel := context.WithTimeout(context.Background(), sqlTimeout) // TODO: handle proper cancellation if bot shuts down
			defer cancel()
			_, err := dl.sqlSession.Db.ExecContext(ctx, "INSERT INTO chat_commands (guild_id, command_key, command_value) VALUES (?, ?, ?);", m.GuildID, strings.ToLower(matches[1]), matches[2])
			if err != nil {
				logrus.WithError(err).
					WithField("key", matches[1]).
					WithField("value", matches[2]).
					Error("Error adding reaction to database.")
				return
			}
			_, err = s.ChannelMessageSend(m.ChannelID, "Added reaction for !"+matches[1]+".")
			if err != nil {
				logrus.WithError(err).
					WithField("key", matches[1]).
					WithField("value", matches[2]).
					WithField("channelId", m.ChannelID).
					Error("Error sending 'added reaction' message to discord channel.")
			}
			logrus.
				WithField("key", matches[1]).
				WithField("value", matches[2]).
				Info("Added reaction.")
			return
		}

		re = regexp.MustCompile(`^[ ]*!([A-Za-z0-9]+)[ ]*$`)
		matches = re.FindStringSubmatch(m.Content)
		if len(matches) == 2 {
			ctx, cancel := context.WithTimeout(context.Background(), sqlTimeout) // TODO: handle proper cancellation if bot shuts down
			defer cancel()
			rows, err := dl.sqlSession.Db.QueryContext(ctx, "SELECT command_key, command_value FROM chat_commands WHERE guild_id = ? AND command_key = ?;", m.GuildID, strings.ToLower(matches[1]))
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					logrus.WithField("key", matches[1]).Info("No reaction found for this key.")
					return
				}
				logrus.WithError(err).
					WithField("key", matches[1]).
					Error("Error fetching reactions from database.")
				return
			}
			if !rows.Next() {
				if rows.Err() != nil {
					logrus.WithError(rows.Err()).
						WithField("key", matches[1]).
						Error("Error fetching reactions from database (Next() returned false).")
				} else {
					logrus.WithField("key", matches[1]).Info("No reaction found for this key.")
				}
				return
			}
			var commandKey string
			var commandValue string
			err = rows.Scan(&commandKey, &commandValue)
			if err != nil {
				logrus.WithError(err).
					WithField("key", matches[1]).
					Error("Error parsing reaction command key and value from from database.")
				return
			}
			_, err = s.ChannelMessageSend(m.ChannelID, commandValue)
			if err != nil {
				logrus.WithError(err).
					WithField("keyFromMessage", matches[1]).
					WithField("key", commandKey).
					WithField("value", commandValue).
					WithField("channelId", m.ChannelID).
					Error("Error sending reaction message to discord channel.")
			}
			return
		}
	}
}
