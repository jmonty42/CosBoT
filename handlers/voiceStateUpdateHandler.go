package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jmonty42/cosbot/errors"
	"github.com/jmonty42/cosbot/types"
)

var cfg *types.Config

// ah, this is what updates when a user switches channels
func AddVoiceStateUpdateHandler(session *discordgo.Session, config *types.Config) {
	cfg = config
	session.AddHandler(func(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
		fmt.Printf("Received VoiceStateUpdate event for user %s\n", event.UserID)

		// we'll need a map: [userId]->{guildId, channelId}
		// initialize that map first with all of the voice states on the server, then
		// add every member for which there is not a voice state (members that are not in a voice
		// channel on the server), where the channelId is just the empty string
		// then when we get this update event, there will be three scenarios
		//
		// old state		new state		scenario
		// ""				*channelId*		joined channel
		// *channelId1*		*channelId2*	moved channel
		// *channelId*		""				left channel
		//
		// additional use case if this bot runs on multiple servers with shared members:
		// switch voice channels to different server (guildids of old vs new states don't match)
		//   - leave message on old server, join message on new server
		//   turns out in this case the bot gets two update events: one from the server the user
		//   left, and one from the server the user joins
		//
		// I am assuming a VoiceStateUpdate event happens when a user joins (new voicestate) and
		// leaves completely (voicestate deleted?), this will require some testing
		oldState := cfg.CachedVoiceStates[event.UserID].VoiceState
		cfg.CachedVoiceStates[event.UserID].VoiceState = event.VoiceState
		fmt.Printf("username: %s\n", cfg.CachedVoiceStates[event.UserID].UserName)
		if oldState != nil && oldState.ChannelID != "" {
			if event.ChannelID != "" {
				// user moved channels
				fmt.Printf("previous channel: %s\n",
					cfg.CachedGuilds[event.GuildID].ChannelNames[oldState.ChannelID])
				fmt.Printf("new channel: %s\n",
					cfg.CachedGuilds[event.GuildID].ChannelNames[event.ChannelID])
				sendNewChannelMessage(session,
					cfg.CachedGuilds[event.GuildID].DefaultMessageChannelId,
					cfg.CachedGuilds[event.GuildID].ChannelNames[event.ChannelID],
					cfg.CachedVoiceStates[event.UserID].UserName,
					false)
			} else {
				// user left
				fmt.Printf("user is no longer in any voice channel on this guild\n")
				sendDisconnectedMessage(session,
					cfg.CachedGuilds[event.GuildID].DefaultMessageChannelId,
					cfg.CachedVoiceStates[event.UserID].UserName)
			}
		} else {
			// user joined the channel
			fmt.Printf("user was not in a channel previously\n")
			fmt.Printf("new channel: %s\n",
				cfg.CachedGuilds[event.GuildID].ChannelNames[event.ChannelID])
			sendNewChannelMessage(session,
				cfg.CachedGuilds[event.GuildID].DefaultMessageChannelId,
				cfg.CachedGuilds[event.GuildID].ChannelNames[event.ChannelID],
				cfg.CachedVoiceStates[event.UserID].UserName,
				true)
		}

	})
}

func sendNewChannelMessage(session *discordgo.Session, messageChannelId string,
	joinedChannelName string, userName string, isJoin bool) {
	var verb string
	if isJoin {
		verb = " joined channel "
	} else {
		verb = " moved to channel "
	}
	message, err := session.ChannelMessageSendTTS(
		messageChannelId, userName+verb+joinedChannelName)
	errors.PanicOnError(err)
	err = session.ChannelMessageDelete(message.ChannelID, message.ID)
	errors.PanicOnError(err)
}

func sendDisconnectedMessage(session *discordgo.Session, messageChannelId string,
	userName string) {
	message, err := session.ChannelMessageSendTTS(messageChannelId,
		userName+" disconnected")
	errors.PanicOnError(err)
	err = session.ChannelMessageDelete(message.ChannelID, message.ID)
	errors.PanicOnError(err)
}
