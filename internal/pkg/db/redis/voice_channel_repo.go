package redis

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kwhk/go-discord-bot/config"
	repo "github.com/kwhk/go-discord-bot/internal/pkg/repository"
)

type RedisVoiceChannelRepo struct {
	client *redis.Client
}

func NewRedisVoiceChannelRepo(client *redis.Client) *RedisVoiceChannelRepo {
	return &RedisVoiceChannelRepo{
		client: config.Redis,
	}
}

func (r RedisVoiceChannelRepo) GetVCOpenSince(ctx context.Context, guildId string) (repo.VoiceChannelOpenSince, error) {
	val, err := r.client.HGetAll(ctx, newKey(guildId, voiceChannelOpenSinceKey)).Result()
	if err == redis.Nil {
		return make(repo.VoiceChannelOpenSince), err
	} else if err != nil {
		log.Printf("Error: Failed to get %s\n", newKey(guildId, voiceChannelOpenSinceKey))
		return nil, err
	}

	var openSince repo.VoiceChannelOpenSince = make(repo.VoiceChannelOpenSince)

	// Convert string to time.Time
	for channelID, v := range val {
		layout := "2006-01-02 15:04:05.999999999 -0700 MST"
		time, err := time.Parse(layout, v)
		if err != nil {
			log.Printf("Error: Failed to parse duration %s\n", val)
			return nil, err
		}

		openSince[channelID] = time
	}

	return openSince, nil
}

func (r RedisVoiceChannelRepo) SetVCOpenSince(ctx context.Context, guildId string, openSince repo.VoiceChannelOpenSince) error {
	var newOpenSince map[string]string = make(map[string]string)
	for channelID, time := range openSince {
		newOpenSince[channelID] = strings.Split(time.String(), " m=")[0]
	}
	_, err := r.client.HSet(ctx, newKey(guildId, voiceChannelOpenSinceKey), newOpenSince).Result()
	if err != nil {
		log.Printf("Error: Failed to set %s\n", newKey(guildId, voiceChannelOpenSinceKey))
		return err
	}

	return nil
}
