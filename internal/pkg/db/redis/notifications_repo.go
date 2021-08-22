package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/kwhk/go-discord-bot/config"
)

type RedisNotificationsRepo struct {
	client *redis.Client
}

func NewNotificationsRepo(client *redis.Client) *RedisNotificationsRepo {
	return &RedisNotificationsRepo{
		client: config.Redis,
	}
}

func (r RedisNotificationsRepo) GetNotificationChannels(ctx context.Context, guildId string) ([]string, error) {
	val, err := r.client.SMembers(ctx, newKey(guildId, notificationChannelsKey)).Result()
	if err == redis.Nil {
		return make([]string, 0), nil
	} else if err != nil {
		log.Printf("Error Failed to get %s\n", newKey(guildId, notificationChannelsKey))
		return nil, err
	}

	return val, nil
}

func (r RedisNotificationsRepo) SetNotificationChannel(ctx context.Context, guildId string, channelId string) error {
	_, err := r.client.SAdd(ctx, newKey(guildId, notificationChannelsKey), channelId).Result()
	if err != nil {
		log.Printf("Error: Failed to set %s\n", newKey(guildId, notificationChannelsKey))
		return err
	}

	return nil
}

func (r RedisNotificationsRepo) RemoveNotificationChannel(ctx context.Context, guildId string, channelId string) error {
	_, err := r.client.SRem(ctx, newKey(guildId, notificationChannelsKey), channelId).Result()
	if err != nil {
		log.Printf("Error: Failed to remove %s\n", newKey(guildId, notificationChannelsKey))
		return err
	}

	return nil
}
