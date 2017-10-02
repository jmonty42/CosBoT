package main

import (
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

	for index, guild := range guilds {
		fmt.Println("Guild #", index, ": ", guild.Name, " (ID: ", guild.ID, ")")
	}

	<-make(chan struct{})
}