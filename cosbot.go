package main

import (
	"fmt"
	"io/ioutil"
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

	session, err := discordgo.New(string(token))
	panicOnError(err)

	fmt.Println(session.Token)

	guilds, err := session.UserGuilds(0, "", "")
	panicOnError(err)

	for index, guild := range guilds {
		fmt.Println("Guild #", index, ": ", guild.Name, " (ID: ", guild.ID, ")")
	}
}