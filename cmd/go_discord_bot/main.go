package main

import (
	"fmt"
	"os"
	"syscall"
	"os/signal"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/kwhk/go-discord-bot/internal/pkg/commands"
)

const (
	GLOBAL_COMMAND = ""
)

func startSession() (*discordgo.Session, error) {
	s, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Invalid bot parameter: %v", err)
		return nil, err
	}
	
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commands.CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Bot is now running. Press CTRL-C to exit.")
	})

	err = s.Open()
	if err != nil {
		fmt.Println("Error opening connection, ", err)
		return nil, err
	}

	return s, nil
}

func main() {
	session, err := startSession()
	if err != nil {
		log.Fatal("Error starting Discord session")
		return
	}

	for _, v := range commands.Commands {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, GLOBAL_COMMAND, v)
		if err != nil {
			fmt.Printf("Cannot create '%v' command: %v\n", v.Name, err)
		}
	}

	defer session.Close()

	// Gracefully close down Discord session.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}