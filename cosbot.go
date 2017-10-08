package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jmonty42/cosbot/errors"
	"github.com/jmonty42/cosbot/handlers"
	"github.com/jmonty42/cosbot/types"
)

const tokenFileName = "TOKEN"
const defaultChannelForMessages = "channel-move-messages"

func main() {
	cfg := types.Config{}

	tokenBytes, err := ioutil.ReadFile(tokenFileName)
	errors.PanicOnError(err)
	cfg.Token = "Bot " + strings.TrimSpace(string(tokenBytes))

	session, err := discordgo.New(cfg.Token)
	errors.PanicOnError(err)

	fmt.Println(session.Token)

	err = session.Open()
	defer session.Close()
	errors.PanicOnError(err)

	userGuilds, err := session.UserGuilds(0, "", "")
	errors.PanicOnError(err)

	cfg.CachedGuilds = make(map[string]*types.CachedGuild)
	cfg.CachedVoiceStates = make(map[string]*types.CachedVoiceState)

	// TODO - break this out as the guildSetup method
	fmt.Printf("Found %d guilds:\n", len(userGuilds))
	for index, userGuild := range userGuilds {
		guild, err := session.Guild(userGuild.ID)
		//put a copy of the current guild object in the cached guilds
		cfg.CachedGuilds[guild.ID] = &types.CachedGuild{
			Guild: guild,
			DefaultMessageChannelId: "",
			ChannelNames:            make(map[string]string)}
		errors.PanicOnError(err)
		fmt.Printf("Guild #%d: %s (ID: %s)\n", index, guild.Name, guild.ID)

		//set the default message channel
		guildChannels, err := session.GuildChannels(guild.ID)
		errors.PanicOnError(err)
		fmt.Printf("  Found %d channels:\n", len(guildChannels))
		for _, channel := range guildChannels {
			cfg.CachedGuilds[guild.ID].ChannelNames[channel.ID] = channel.Name
			var channelType string
			switch channel.Type {
			case discordgo.ChannelTypeGuildText:
				channelType = "Text"
			case discordgo.ChannelTypeGuildVoice:
				channelType = "Voice"
			case discordgo.ChannelTypeGuildCategory:
				channelType = "Category"
			default:
				channelType = "unknown"
			}
			fmt.Printf("  %s (ID: %s) %s\n", channel.Name, channel.ID, channelType)

			// sets the first text channel it finds as the default message channel
			// for that server, then if it finds another text channel that matches the default
			// channel name it sets the message channel to that one
			if (cfg.CachedGuilds[guild.ID].DefaultMessageChannelId == "" ||
				channel.Name == defaultChannelForMessages) &&
				channel.Type == discordgo.ChannelTypeGuildText {
				cfg.CachedGuilds[guild.ID].DefaultMessageChannelId = channel.ID
			}
		}

		//cache the guild's voice states
		fmt.Printf("  Found %d 'VoiceStates':\n", len(guild.VoiceStates))
		for _, voiceState := range guild.VoiceStates {
			user, err := session.User(voiceState.UserID)
			if err != nil {
				fmt.Printf("Got error when trying to get user with id: %s : %s\n",
					voiceState.UserID, err.Error())
			}
			channel, err := session.Channel(voiceState.ChannelID)
			if err != nil {
				fmt.Printf("Got error when trying to get channel with id: %s : %s\n",
					voiceState.ChannelID, err.Error())
			}
			fmt.Printf("    %s: %s\n", user.Username, channel.Name)
			cfg.CachedVoiceStates[voiceState.UserID] = &types.CachedVoiceState{
				VoiceState: voiceState, UserName: user.Username}
		}
		//iterate over the guild members, for those that didn't have corresponding voice states,
		//add an empty voice state to the cache
		for _, member := range guild.Members {
			if _, present := cfg.CachedVoiceStates[member.User.ID]; !present {
				cfg.CachedVoiceStates[member.User.ID] = &types.CachedVoiceState{
					VoiceState: nil, UserName: member.User.Username}
			}
		}

		channel, err := session.Channel(cfg.CachedGuilds[guild.ID].DefaultMessageChannelId)
		if err != nil {
			fmt.Printf("Got error when trying to get the channel %s: %s",
				cfg.CachedGuilds[guild.ID].DefaultMessageChannelId, err.Error())
		} else {
			fmt.Println(" Default channel: ", channel.Name)
		}
	}

	handlers.AddGuildUpdateHandler(session)
	handlers.AddVoiceStateUpdateHandler(session, &cfg)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Bot running. Press 'Enter' to quit ...")
	reader.ReadBytes('\n')
}

// haven't actually tested this out yet, but will be needed for the command to set the
// channel that messages should be sent to
func getChannelIdFromName(session *discordgo.Session, channelName string,
	guildId string) (string, error) {
	channels, err := session.GuildChannels(guildId)
	if err != nil {
		return "", err
	}

	for _, channel := range channels {
		if channel.Name == channelName && channel.Type == discordgo.ChannelTypeGuildText {
			return channel.ID, nil
		}
	}

	return "", fmt.Errorf("Channel not found: '%s' on server: '%s'", channelName, guildId)
}
