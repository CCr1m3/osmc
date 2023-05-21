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
			log.Fatal("failed to create channel group BotChannels: ", err.Error())
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
	return "Hello, my name is <something> and I'm one of this Discord's assistants. I'm here to assign you ranks based on your performance in the game.\n\nUse `/assignRank` so I can connect the given Omega Strikers and assign a rank to you! Do mind that connecting an account that does not belong to you is punishable.\nAlso, do mind that you will be assigned to Rookie if you are outside the top 10,000 on the global leaderboard yet.\nIn case you wish to disconnect your Omega Strikers account, please contact a moderator.\n\nAnd if you have improved your rank after I have already assigned a rank to you, use `/updateRank` so I can update your rank so everyone can see your progress!\n\nGood luck in your journey to Omega!"
}
