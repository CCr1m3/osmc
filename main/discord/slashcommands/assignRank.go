package slashcommands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ccr1m3/osmc/main/prometheus"
	"github.com/ccr1m3/osmc/main/static"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type AssignRank struct{}

func (p AssignRank) Name() string {
	return "sync"
}

func (p AssignRank) Description() string {
	return "Allows you to sync to an omega strikers account."
}

func (p AssignRank) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionSendMessages)
	return &perm
}

func (p AssignRank) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:        "Username",
			Description: "Username in omega strikers",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	}
}

func (p AssignRank) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	playerID := i.Member.User.ID
	username := strings.ToLower(optionMap["Username"].StringValue())
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	log.WithFields(log.Fields{
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): playerID,
		string(static.UsernameKey): username,
	}).Info("AssignRank slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "AssignRank slashcommand invoked. Please wait...",
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

	err = prometheus.LinkPlayerToUsername(ctx, playerID, username)
	if err != nil {
		log.Errorf("failed to sync player %s with username %s: "+err.Error(), playerID, username)
		switch {
		case errors.Is(err, static.ErrUsernameNotFound):
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.UsernameKey): username,
			}).Warning("failed to sync player, username invalid")
			message = fmt.Sprintf("Could not find username: %s", username)
		case errors.Is(err, static.ErrUserAlreadyLinked):
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
			}).Warning("failed to link player, user already linked")
			message = "You have already linked an omega strikers account. Please contact a mod if you want to unsync."
		case errors.Is(err, static.ErrUsernameAlreadyLinked):
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.UsernameKey): username,
			}).Warning("failed to sync player, username already synced")
			message = fmt.Sprintf("%s is already linked to an account. Please contact a mod if you think you are the rightful owner of the account.", username)
		default:
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.UsernameKey): username,
				string(static.ErrorKey):    err.Error(),
			}).Error("failed to sync player")
			message = fmt.Sprintf("Failed to link to %s.", username)
		}
		return
	}
	log.WithFields(log.Fields{
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): i.Member.User.ID,
		string(static.UsernameKey): username,
	}).Info("player successfully synced")
	message = fmt.Sprintf("Successfully synchronized to %s.", username)
}

type UnassignRank struct{}

func (p UnassignRank) Name() string {
	return "unsync"
}

func (p UnassignRank) Description() string {
	return "Allow mods to unsync someone from his Omega Strikers' account."
}

func (p UnassignRank) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionModerateMembers)
	return &perm
}

func (p UnassignRank) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:        "Discord user",
			Description: "User in Discord",
			Type:        discordgo.ApplicationCommandOptionUser,
			Required:    true,
		},
	}
}

func (p UnassignRank) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	callerID := i.Member.User.ID
	user := optionMap["Discord user"].UserValue(s)
	log.WithFields(log.Fields{
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): callerID,
		string(static.PlayerIDKey): user.ID,
	}).Info("UnassignRank slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "UnassignRank slash command invoked. Please wait...",
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
	if i.Member.Permissions&discordgo.PermissionModerateMembers != discordgo.PermissionModerateMembers {
		message = "You do not have the permission to unsync."
		return
	}
	err = prometheus.UnlinkPlayer(ctx, user.ID)
	if err != nil {
		if errors.Is(err, static.ErrUserNotLinked) {
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.PlayerIDKey): user.ID,
			}).Warning("player is not synced")
			message = fmt.Sprintf("%s has not synchronized an omega strikers account.", user.Mention())
		} else {
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.PlayerIDKey): user.ID,
				string(static.ErrorKey):    err.Error(),
			}).Error("failed to unsync player")
			message = fmt.Sprintf("Failed to unsync %s.", user.Mention())
		}
		return
	}
	message = fmt.Sprintf("Successfully unsynced %s.", user.Mention())
}
