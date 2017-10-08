package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// this doesn't update when a user switches voice channels, not sure if we'll need this
func AddGuildUpdateHandler(session *discordgo.Session) {
	session.AddHandler(func(session *discordgo.Session, event *discordgo.GuildUpdate) {
		fmt.Printf("Received GuildUpdate event for guild: %s\n", event.ID)
	})
}
