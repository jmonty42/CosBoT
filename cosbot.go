package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type CachedGuild struct {
	guild                   *discordgo.Guild
	defaultMessageChannelId string //the channel that messages will be sent to
}

type Config struct {
	token string
	// I'm sticking this in the 'config' because it's part of the state that should be
	// saved when the bot shuts down, because the channel that the messages get sent to
	// can be changed at runtime with a command (eventually)
	cachedGuilds map[string]*CachedGuild //key=guildId, value=CachedGuild object
}

const tokenFileName = "TOKEN"
const defaultChannelForMessages = "test-channel"

func panicOnError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	cfg := Config{}

	tokenBytes, err := ioutil.ReadFile(tokenFileName)
	panicOnError(err)
	cfg.token = "Bot " + strings.TrimSpace(string(tokenBytes))

	session, err := discordgo.New(cfg.token)
	panicOnError(err)

	fmt.Println(session.Token)

	err = session.Open()
	defer session.Close()
	panicOnError(err)

	userGuilds, err := session.UserGuilds(0, "", "")
	panicOnError(err)

	cfg.cachedGuilds = make(map[string]*CachedGuild)

	// TODO - break this out as the findDefaultChannels method
	fmt.Printf("Found %d guilds:\n", len(userGuilds))
	for index, userGuild := range userGuilds {
		guild, err := session.Guild(userGuild.ID)
		cfg.cachedGuilds[guild.ID] = &CachedGuild{guild, ""}
		panicOnError(err)
		fmt.Printf("Guild #%d: %s (ID: %s)\n", index, guild.Name, guild.ID)
		guildChannels, err := session.GuildChannels(guild.ID)
		panicOnError(err)
		fmt.Printf("  Found %d channels:\n", len(guildChannels))
		for _, channel := range guildChannels {
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

			// I thought this would give a list of users in each channel ... it does not
			// leaving it here for posterity (will remove in a bit)
			// fmt.Printf("    Recipients: (%d)\n", len(channel.Recipients))
			// for _, user := range channel.Recipients {
			// 	fmt.Println("      ", user.Username)
			// }

			// sets the first text channel it finds as the default message channel
			// for that server, then if it finds another text channel that matches the default
			// channel name it sets the message channel to that one
			if (cfg.cachedGuilds[guild.ID].defaultMessageChannelId == "" ||
				channel.Name == defaultChannelForMessages) &&
				channel.Type == discordgo.ChannelTypeGuildText {
				cfg.cachedGuilds[guild.ID].defaultMessageChannelId = channel.ID
			}
		}

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
		}

		channel, err := session.Channel(cfg.cachedGuilds[guild.ID].defaultMessageChannelId)
		if err != nil {
			fmt.Printf("Got error when trying to get the channel %s: %s",
				cfg.cachedGuilds[guild.ID].defaultMessageChannelId, err.Error())
		} else {
			fmt.Println(" Default channel: ", channel.Name)
		}
	}

	addGuildUpdateHandler(session)
	addVoiceStateUpdateHandler(session)

	// TODO - remove this as it's just code to learn how to send messages
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the channel ID of a channel to send a text message to: ")
	channelId, err := reader.ReadString('\n')
	panicOnError(err)

	if strings.TrimSpace(channelId) != "" {
		// sends a normal message to a text channel, will not be TTS unless the users have it set
		// for that channel
		_, err = session.ChannelMessageSend(strings.TrimSpace(channelId), "Test text message")
		panicOnError(err)
	}

	// TODO - remove this as it's just code to learn how to send messages
	fmt.Println("Enter the channel ID of a channel to send a TTS message to: ")
	channelId, err = reader.ReadString('\n')
	panicOnError(err)

	if strings.TrimSpace(channelId) != "" {
		// sends a message to a text channel that will be read out as TTS to any users monitoring
		// that channel (as far as I can tell)
		message, err := session.ChannelMessageSendTTS(
			strings.TrimSpace(channelId), "Test TTS message (a different one)")
		panicOnError(err)
		// tested this out, and it works! deleting the TTS message right after sending it
		// will read it out, but not clutter up the channel (could be a configurable option)
		err = session.ChannelMessageDelete(message.ChannelID, message.ID)
		if err != nil {
			fmt.Println("Got an error when trying to delete the message: ", err.Error())
		}
	}

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

// either this doesn't update when a user switches voice channels, or the input
// blocks above
func addGuildUpdateHandler(session *discordgo.Session) {
	session.AddHandler(func(session *discordgo.Session, event *discordgo.GuildUpdate) {
		fmt.Printf("Received GuildUpdate event for guild: %s\n", event.ID)
	})
}

func addVoiceStateUpdateHandler(session *discordgo.Session) {
	session.AddHandler(func(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
		fmt.Printf("Received VoiceStateUpdate event for user %s", event.UserID)
	})
}
