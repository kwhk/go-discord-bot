package monitor

import (
	"context"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kwhk/go-discord-bot/config"
	"github.com/kwhk/go-discord-bot/internal/pkg/db/redis"
)

type Config struct {
	Session  *discordgo.Session
	Interval time.Duration
}

type Metrics struct {
	*Config
	Guilds        *Guilds
	VoiceChannels *VoiceChannels
}

func NewMetrics(config *Config) *Metrics {
	metrics := &Metrics{
		Config: config,
	}

	metrics.newGuilds()
	metrics.newVoiceChannels()

	return metrics
}

func (metrics *Metrics) Monitor(ctx context.Context) {
	go metrics.Guilds.Monitor(ctx)
	go metrics.VoiceChannels.Monitor(ctx)
}

func (metrics *Metrics) newGuilds() {
	metrics.Guilds = &Guilds{
		Session:  metrics.Session,
		Interval: metrics.Interval,
		Cache:    &GuildsCache{Mutex: &sync.Mutex{}},
	}
}

func (metrics *Metrics) newVoiceChannels() {
	metrics.VoiceChannels = &VoiceChannels{
		Session:           metrics.Session,
		Interval:          metrics.Interval,
		VoiceChannelRepo:  redis.NewRedisVoiceChannelRepo(config.Redis),
		NotificationsRepo: redis.NewNotificationsRepo(config.Redis),
	}
}
