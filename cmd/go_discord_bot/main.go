package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/kwhk/go-discord-bot/config"
	"github.com/kwhk/go-discord-bot/internal/pkg/callbacks"
	"github.com/kwhk/go-discord-bot/internal/pkg/commands"
	"github.com/kwhk/go-discord-bot/internal/pkg/monitor"
)

const (
	globalCommand   = ""
	monitorInterval = 5 * time.Minute
)

func startSession(ctx context.Context) (*discordgo.Session, error) {
	s, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Invalid bot parameter: %v", err)
		return nil, err
	}

	metrics := monitor.NewMetrics(&monitor.Config{
		Session:  s,
		Interval: monitorInterval,
	})

	addCallbacks(s)

	err = s.Open()
	if err != nil {
		fmt.Println("Error opening connection, ", err)
		return nil, err
	}

	metrics.Monitor(ctx)

	return s, nil
}

func addCallbacks(session *discordgo.Session) {
	session.AddHandler(callbacks.InitCommands)
	session.AddHandler(callbacks.Ready)
}

func main() {
	monitorCtx, cancelMonitorCtx := context.WithCancel(context.Background())
	defer cancelMonitorCtx()

	session, err := startSession(monitorCtx)
	if err != nil {
		log.Fatal("Error starting Discord session")
		return
	}

	for _, v := range commands.Commands {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, globalCommand, v)
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
