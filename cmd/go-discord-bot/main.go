package main

import (
	"fmt"
	"os"
	"syscall"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/kwhk/go-discord-bot/config"
	"github.com/kwhk/go-discord-bot/pkg/commands"
)

func main() {
	dg, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("Error creating discord session: ", err)
		return
	}

	dg.AddHandler(messageSender)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection, ", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	// Gracefully close down Discord session.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageSender(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Bot should not reply to itself.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content[0] == config.Prefix {
		go commands.Controller(s, m, m.Content[1:])
	}
}