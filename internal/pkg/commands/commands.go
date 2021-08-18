package commands

import (
	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Ping bot",
	},
	{
		Name:        "options",
		Description: "Command for demonstrating options",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "string-option",
				Description: "String option",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "integer-option",
				Description: "Integer option",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "bool-option",
				Description: "Boolean option",
				Required:    true,
			},

			// Required options must be listed first since optional parameters
			// always come after when they're used.
			// The same concept applies to Discord's Slash-commands API

			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel-option",
				Description: "Channel option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user-option",
				Description: "User option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role-option",
				Description: "Role option",
				Required:    false,
			},
		},
	},
}

var CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping":    ping,
	"options": options,
}
