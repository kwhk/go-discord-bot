// Package monitor provides background monitoring for entities in Discord that are tracked periodically.
package monitor

import (
	"fmt"
	"context"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Guilds struct {
	Session *discordgo.Session
	Interval time.Duration
	Cache *GuildsCache
}

type GuildsCache struct {
	Mutex *sync.Mutex
	guildList []*discordgo.Guild
	numGuilds int
}

func (guilds *Guilds) Monitor(ctx context.Context) {
	updateTicker := time.NewTicker(guilds.Interval)
	defer updateTicker.Stop()

	for {
		select {
		case <-updateTicker.C:
			guilds.update()
		case <-ctx.Done():
			return
		}
	}
}

// update monitors which guilds the bot is connected to and logs recent connection/disconnections.
func (guilds *Guilds) update() {
	guilds.Cache.Mutex.Lock()
	defer guilds.Cache.Mutex.Unlock()

	originalCount := guilds.Cache.numGuilds
	newCount := len(guilds.Session.State.Guilds)

	switch {
	case newCount == originalCount:
		return
	case newCount > originalCount && originalCount != 0:
		botName := guilds.Session.State.User.Username
		newGuild := guilds.Session.State.Guilds[newCount-1]
		fmt.Printf("%s joined new guild, '%s'\n", botName, newGuild.Name)
	case newCount < originalCount:
		botName := guilds.Session.State.User.Username
		fmt.Printf("%s removed from guild\n", botName)
	}

	guilds.Cache.numGuilds = newCount
	guilds.Cache.guildList = guilds.Session.State.Guilds
}