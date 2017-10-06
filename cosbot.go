package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	token string
	// I'm sticking this in the 'config' because it's part of the state that should be
	// saved when the bot shuts down, because the channel that the messages get sent to
	// can be changed at runtime with a command (eventually)
	messageChannels map[string]string //key=guildId, value=default channelId for messages
}

const tokenFileName = "TOKEN"
const defaultChannelForMessages = "general"

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

	guilds, err := session.UserGuilds(0, "", "")
	panicOnError(err)

	cfg.messageChannels = make(map[string]string)

	// TODO - break this out as the findDefaultChannels method
	fmt.Println("Found ", len(guilds), " guilds:")
	for index, guild := range guilds {
		fmt.Println("Guild #", index, ": ", guild.Name, " (ID: ", guild.ID, ")")
		guildChannels, err := session.GuildChannels(guild.ID)
		panicOnError(err)
		fmt.Println("  Found ", len(guildChannels), " channels:")
		for index, channel := range guildChannels {
			fmt.Println("  ", channel.Name, " (ID: ", channel.ID,
				") Text: ", channel.Type == discordgo.ChannelTypeGuildText)
			// sets the first text channel it finds as the default message channel
			// for that server, then if it finds another text channel that matches the default
			// channel name it sets the message channel to that one
			if (index == 0 || channel.Name == defaultChannelForMessages) &&
				channel.Type == discordgo.ChannelTypeGuildText {
				cfg.messageChannels[guild.ID] = channel.ID
			}
		}
		channel, _ := session.Channel(cfg.messageChannels[guild.ID])
		fmt.Println(" Default channel: ", channel.Name)
	}

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
