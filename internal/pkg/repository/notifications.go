package repository

import (
	"context"
)

type NotificationsRepo interface {
	// Get channels to send notifications to
	GetNotificationChannels(ctx context.Context, guildId string) ([]string, error)
	// Set channels to send notifications to and return success if channel hasn't been added before
	SetNotificationChannel(ctx context.Context, guildId string, channelId string) error
	RemoveNotificationChannel(ctx context.Context, guildId string, channelId string) error
}
