package slashcommands

import (
	"context"
	"errors"
	"fmt"

	"github.com/ccr1m3/osmc/main/db"
	"github.com/ccr1m3/osmc/main/prometheus"
	"github.com/ccr1m3/osmc/main/static"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type UpdateRank struct{}

func (p UpdateRank) Name() string {
	return "updaterank"
}

func (p UpdateRank) Description() string {
	return "Allows you to update your Discord role using your connected Omega Strikers account."
}

func (p UpdateRank) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionSendMessages)
	return &perm
}

func (p UpdateRank) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (p UpdateRank) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	playerID := i.Member.User.ID
	log.WithFields(log.Fields{
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): playerID,
	}).Info("updateRank slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "UpdateRank slash command invoked. Please wait...",
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to send message")
		return
	}
	var message string
	defer func() {
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &message,
		})
		if err != nil {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to edit message")
		}
	}()

	err = prometheus.UpdateRankIfNeeded(ctx, playerID)
	if err != nil {
		switch {
		case errors.Is(err, static.ErrRankUpdateTooFast):
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
			}).Warning("player update too fast")
			message = "You have updated your account recently. Please wait before using this command again."
		case errors.Is(err, static.ErrUserNotLinked):
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
			}).Warning("player is not connected")
			message = "You are not connected to an Omega Strikers account. Please use \"assignRank\" first."
		default:
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.ErrorKey):    err.Error(),
			}).Error("failed to update player")
			message = "Failed to update your rank."
		}
		return
	}
	account, err := db.GetPlayerByID(ctx, playerID)
	if err != nil {
		log.Warning("failed to get player", err.Error())
		message = "Successfully updated your rank."
		return
	}
	rank := prometheus.GetRank(account.Elo).Name
	message = fmt.Sprintf("Successfully updated your rank to %s.", rank)
}
