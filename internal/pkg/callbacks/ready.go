package callbacks

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func Ready(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println("Bot is now running. Press CTRL-C to exit.")	
}