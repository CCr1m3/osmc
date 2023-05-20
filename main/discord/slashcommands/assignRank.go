package slashcommands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ccr1m3/osmc/main/db"
	"github.com/ccr1m3/osmc/main/prometheus"
	"github.com/ccr1m3/osmc/main/static"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type AssignRank struct{}

func (p AssignRank) Name() string {
	return "assignrank"
}

func (p AssignRank) Description() string {
	return "Assign yourself a rank based on your Omega Strikers rank."
}

func (p AssignRank) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionSendMessages)
	return &perm
}

func (p AssignRank) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:        "username",
			Description: "Username in Omega Strikers",
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
	username := strings.ToLower(optionMap["username"].StringValue())
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

	account, err := prometheus.LinkPlayerToUsername(ctx, playerID, username)
	if err != nil {
		log.Errorf("failed to sync player %s with username %s: "+err.Error(), playerID, username)
		switch {
		case errors.Is(err, static.ErrUsernameNotFound):
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.UsernameKey): username,
			}).Warning("failed to connect player, username invalid")
			message = fmt.Sprintf("Could not find username: %s", username)
		case errors.Is(err, static.ErrUsernameAlreadyLinked):
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.UsernameKey): username,
			}).Warning("failed to sync player, username already synced")
			message = fmt.Sprintf("%s is already connected to another user. Please contact a moderator if you think you are the rightful owner of the account.", username)
		case errors.Is(err, static.ErrUserAlreadyLinked):
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
			}).Warning("failed to connect player, user already connect")
			message = fmt.Sprintf("You are already connected to %s. Please contact a moderator if you wish to disconnect from this account.", account.PlayerUsername)
		default:
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.UsernameKey): username,
				string(static.ErrorKey):    err.Error(),
			}).Error("failed to connect player")
			message = fmt.Sprintf("Failed to connect to %s.", username)
		}
		return
	}
	account, err = db.GetPlayerByID(ctx, playerID)
	if err != nil {
		log.Warning("failed to get player", err.Error())
		message = fmt.Sprintf("Successfully assigned rank to %s.", i.Member.User.Mention())
		return
	}
	rank := prometheus.GetRank(account.Elo).Name
	log.WithFields(log.Fields{
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): i.Member.User.ID,
		string(static.UsernameKey): username,
	}).Info("player successfully synced")
	message = fmt.Sprintf("Successfully assigned %s to %s!", rank, i.Member.User.Mention())
}

type UnassignRank struct{}

func (p UnassignRank) Name() string {
	return "unassignrank"
}

func (p UnassignRank) Description() string {
	return "Allows mods to disconnect someone of their Omega Strikers account."
}

func (p UnassignRank) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionModerateMembers)
	return &perm
}

func (p UnassignRank) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:        "discorduser",
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
	user := optionMap["discorduser"].UserValue(s)
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
		message = "You do not have the permission to disconnect someone."
		return
	}
	err = prometheus.UnlinkPlayer(ctx, user.ID)
	if err != nil {
		if errors.Is(err, static.ErrUserNotLinked) {
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.PlayerIDKey): user.ID,
			}).Warning("player is not connected")
			message = fmt.Sprintf("%s is not connected to an Omega Strikers account.", user.Mention())
		} else {
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.PlayerIDKey): user.ID,
				string(static.ErrorKey):    err.Error(),
			}).Error("failed to disconnect player")
			message = fmt.Sprintf("Failed to disconnect %s.", user.Mention())
		}
		return
	}
	message = fmt.Sprintf("Successfully disconnected %s!", user.Mention())
}
