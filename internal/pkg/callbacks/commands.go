package callbacks

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kwhk/go-discord-bot/internal/pkg/commands"
)

func InitCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := commands.CommandHandlers[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}
