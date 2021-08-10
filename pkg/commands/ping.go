package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
	if err != nil {
		fmt.Println("Error sending DM message: ", err)
		s.ChannelMessageSend(m.ChannelID, "Failed to send you a DM.")
	}
}