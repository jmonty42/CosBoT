# CosBoT

Writing a little discord bot to better learn go and provide some functionality that we miss from Mumble, mainly TTS announcements when someone leaves/joins a channel or the server.

You'll need to get the discordgo library to build it locally with ```go get github.com/bwmarrin/discordgo```
(I'll probably end up using a dependency manager like [dep](https://github.com/golang/dep))

Currently, it expects a [discord bot token](https://discordapp.com/developers/docs/intro) in a file named "TOKEN". 

Right now it will send a TTS message to either a channel named 'channel-move-messages' or the first text channel it finds if that channel doesn't exist. It will send a message to that channel each time a user joins or leaves a voice channel on the server.

Added a [changelog!](CHANGELOG.md)

Roadmap:
* Move all the setup code out of main.go into a sub package
* Use a human-readable, parsable config file
  * Something like JSON or YAML
  * Contain token, default channel for messages, anything we end up needing
* Some commands we should support:
  * enable/disable announcements
  * toggle whether announcements go to a specific channel or directly to users
  * set which channel announcements are set to go to
  * users should be able to turn off direct messages if that is what is being used (!optout or something)
  * whitelist certain users to be able to use certain commands
* Any of the settings that are set with the commands above should be saved to a config file so that they are not lost when the bot host is restarted
* Unit tests should be added at some point
* Debug print statements should be replaced with an actual logger