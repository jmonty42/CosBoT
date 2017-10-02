package main

import (
	"os"
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const tokenFileName = "TOKEN"

func panicOnError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	token, err := ioutil.ReadFile(tokenFileName)
	panicOnError(err)

	session, err := discordgo.New("Bot " + strings.TrimSpace(string(token)))
	panicOnError(err)

	fmt.Println(session.Token)

	err = session.Open()
	defer session.Close()
	panicOnError(err)

	guilds, err := session.UserGuilds(0, "", "")
	panicOnError(err)

	fmt.Println("Found ", len(guilds), " guilds:")
	for index, guild := range guilds {
		fmt.Println("Guild #", index, ": ", guild.Name, " (ID: ", guild.ID, ")")
		guildChannels, err := session.GuildChannels(guild.ID)
		panicOnError(err)
		fmt.Println("  Found ", len(guildChannels), " channels:")
		for _, channel := range guildChannels {
			fmt.Println("  ", channel.Name, " (ID: ", channel.ID, ")")
		}
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the channel ID of a channel to send a text message to: ")
	channelId, err := reader.ReadString('\n')
	panicOnError(err)

	// sends a normal message to a text channel, will not be TTS unless the users have it set
	// for that channel
	_, err = session.ChannelMessageSend(strings.TrimSpace(channelId), "Test text message")
	panicOnError(err)

	fmt.Println("Enter the channel ID of a channel to send a TTS message to: ")
	channelId, err = reader.ReadString('\n')
	panicOnError(err)

	// sends a message to a text channel that will be read out as TTS to any users monitoring
	// that channel (as far as I can tell)
	_, err = session.ChannelMessageSendTTS(strings.TrimSpace(channelId), "Test TTS message")

	fmt.Println("Bot running. Press 'Enter' to quit ...")
	reader.ReadBytes('\n')
}