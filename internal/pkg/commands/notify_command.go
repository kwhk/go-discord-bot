package commands

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/kwhk/go-discord-bot/config"
	"github.com/kwhk/go-discord-bot/internal/pkg/db/redis"
)

func notify(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()
	repo := redis.NewNotificationsRepo(config.Redis)

	var response string

	notificationChannels, err := repo.GetNotificationChannels(ctx, i.GuildID)
	if err != nil {
		response = "Bot failed to execute command."
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})
		return
	}

	// Check if channel is already listening.
	var found bool = false
	for _, channelId := range notificationChannels {
		if channelId == i.ChannelID {
			found = true
		}
	}

	// If channel was already listening, disable it.
	if found {
		err := repo.RemoveNotificationChannel(ctx, i.GuildID, i.ChannelID)
		if err != nil {
			response = "Bot failed to execute command."
		} else {
			response = "Bot will no longer post automatic notifications to this channel."
		}
	} else {
		err := repo.SetNotificationChannel(ctx, i.GuildID, i.ChannelID)
		if err != nil {
			response = "Bot failed to execute command."
		} else {
			response = "Bot will now post automatic notifications to this channel."
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}
