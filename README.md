# CosBoT

Writing a little discord bot to better learn go and provide some functionality that we miss from Mumble, mainly TTS announcements when someone leaves/joins a channel or the server.

You'll need to get the discordgo library to build it locally with ```go get github.com/bwmarrin/discordgo```
(I'll probably end up using a dependency manager like [dep](https://github.com/golang/dep))

Currently, it expects a [discord bot token](https://discordapp.com/developers/docs/intro) in a file named "TOKEN". 

Planned features:
* Human-readable, parsable config file
  * Something like JSON or YAML
  * Contain token, default channel for messages, anything we end up needing
* Handler for announcing when someone leaves/joins the server
  * As far as I can tell, we'll have to have a handler for the [GuildUpdate event](https://godoc.org/github.com/bwmarrin/discordgo#GuildUpdate). When we catch this event, we'll have to compare the [Presence](https://godoc.org/github.com/bwmarrin/discordgo#Presence) slice of the new Guild object with the Presence slice of the cached Guild object we'll have to keep around and update on each GuildUpdate event.  
   The Presence object represents a guild members presence (online, offline, in a game, etc). I think that is the only way to see if a user connects or disconnects from the server.  
   **This might not be as useful as announcenments for channel joins/leaves**
* Handler for announcing when someone leaves/joins a channel
  * We'll want to handle [ChannelUpdate events](https://godoc.org/github.com/bwmarrin/discordgo#ChannelUpdate) for this one. The new Channel object's 'recipients' slice should be compared to the cached Channel object's to determine if someone has left or joined.
  * A tricky thing about this: you can't send TTS messages to voice channels. So we'll probably need to figure out whether we just want to announce all channel moves in the general channel, or send private messages to everybody in a particular channel when someone leaves/joins that channel.
  * [Session.ChannelMessageSendTTS](https://godoc.org/github.com/bwmarrin/discordgo#Session.ChannelMessageSendTTS) is what we want to use for these types of messages (unless we go the direct-message route - haven't looked into how to send those)
* Some commands we should support:
  * enable/disable announcements
  * toggle whether announcements go to a specific channel or directly to users
  * set which channel announcements are set to go to
  * users should be able to turn off direct messages if that is what is being used (!optout or something)
* Any of the settings that are set with the commands above should be saved to a config file so that they are not lost when the bot host is restarted