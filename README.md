# CosBoT

Writing a little discord bot to better learn go and provide some functionality that we miss from Mumble, mainly TTS announcements when someone leaves/joins a channel or the server.

You'll need to get the discordgo library to build it locally with ```go get github.com/bwmarrin/discordgo```
(I'll probably end up using a dependency manager like [dep](https://github.com/golang/dep))

Currently, it expects a [discord bot token](https://discordapp.com/developers/docs/intro) in a file named "TOKEN". 