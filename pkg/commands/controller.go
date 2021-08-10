package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

// Map action (key) with commands
var commands = map[string][]string {
	"ping": {"ping", "p"},
}

// parseCommand receives a command and parses it to an action
func parseCommand(command string) string {
	for action, arr := range commands {
		for _, val := range arr {
			if val == command {
				return action
			}
		}
	}

	return "No such command found"
}

func Controller(s *discordgo.Session, m *discordgo.MessageCreate, command string) {
	action := parseCommand(command)
	fmt.Println(action)

	switch action {
	case "ping":
		Ping(s, m)
		fmt.Println("ping was executed.")
	}
}