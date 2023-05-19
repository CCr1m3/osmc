package discord

import (
	"github.com/bwmarrin/discordgo"

	log "github.com/sirupsen/logrus"
)

var BotChannels *discordgo.Channel
var HowToChannel *discordgo.Channel

func initChannels() error {
	channels, err := session.GuildChannels(GuildID)
	if err != nil {
		log.Error("failed to get guild channels: ", err.Error())
	}
	for _, channel := range channels {
		if channel.Name == "BotChannels" {
			BotChannels = channel
		}
		if channel.Name == "rank-roles" {
			HowToChannel = channel
		}
	}
	if BotChannels == nil {
		BotChannels, err = session.GuildChannelCreate(GuildID, "BotChannels", discordgo.ChannelTypeGuildCategory)
		if err != nil {
			log.Fatal("failed to create channel group Ai.Mi: ", err.Error())
		}
	}
	if HowToChannel == nil {
		HowToChannel, err = session.GuildChannelCreateComplex(GuildID, discordgo.GuildChannelCreateData{Name: "rank-roles", Type: discordgo.ChannelTypeGuildText, ParentID: BotChannels.ID})
		if err != nil {
			log.Fatal("failed to create channel how-to: ", err.Error())
		}
		err = session.ChannelPermissionSet(HowToChannel.ID, GuildID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
		if err != nil {
			log.Fatal("failed to lock channel how-to: ", err.Error())
		}
		err = session.ChannelPermissionSet(HowToChannel.ID, session.State.User.ID, discordgo.PermissionOverwriteTypeMember, discordgo.PermissionSendMessages, 0)
		if err != nil {
			log.Fatal("failed to open channel how-to for bot: ", err.Error())
		}
	}
	err = initHowTo()
	if err != nil {
		log.Fatal("failed to init channel how-to: ", err.Error())
	}
	return nil
}

func initHowTo() error {
	channelMessages, err := session.ChannelMessages(HowToChannel.ID, 100, "", "", "")
	if err != nil {
		return err
	}
	if len(channelMessages) == 0 {
		_, err := session.ChannelMessageSend(HowToChannel.ID, howtomessage())
		if err != nil {
			return err
		}
	} else {
		_, err := session.ChannelMessageEdit(HowToChannel.ID, channelMessages[0].ID, howtomessage())
		if err != nil {
			return err
		}
	}
	return nil
}

func howtomessage() string {
	return "lorem ipsum dolor"
}
