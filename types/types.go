package types

import (
	"github.com/bwmarrin/discordgo"
)

type CachedGuild struct {
	Guild                   *discordgo.Guild
	DefaultMessageChannelId string            //the channel that messages will be sent to
	ChannelNames            map[string]string //key=channel id, value=channel name
}

type CachedVoiceState struct {
	VoiceState *discordgo.VoiceState
	UserName   string //cache the username with the voice state so we don't have to query for it
}

type Config struct {
	Token string
	// I'm sticking this in the 'config' because it's part of the state that should be
	// saved when the bot shuts down, because the channel that the messages get sent to
	// can be changed at runtime with a command (eventually)
	CachedGuilds map[string]*CachedGuild //key=guildId, value=CachedGuild object
	// one limitation of this is if this bot runs on multiple servers that share members, those
	// members will have a different VoiceState for each server. Since we only have one server
	// right now where this will run, we won't worry about it
	CachedVoiceStates map[string]*CachedVoiceState //key=userId
}
