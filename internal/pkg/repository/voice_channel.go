/* Abstract away from the specific implementation of DB operations and instead define
its expected behaviour. This makes it easier when changing DB (e.g. from Redis to MongoDB)
in the future. */
package repository

import (
	"context"
	"time"
)

type VoiceChannelOpenSince map[string]time.Time

type VoiceChannelRepo interface {
	// Get all voice channels that have been active
	GetVCOpenSince(ctx context.Context, guildId string) (VoiceChannelOpenSince, error)
	SetVCOpenSince(
		ctx context.Context,
		guildId string,
		time VoiceChannelOpenSince,
	) error
}
