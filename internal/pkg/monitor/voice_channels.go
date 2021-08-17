package monitor

import (
	"fmt"
	"context"
	"time"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kwhk/go-discord-bot/internal/pkg/repository"
)

const (
	timeStampSeparator = ": "
)

type VoiceChannels struct {
	Session *discordgo.Session
	Interval time.Duration
	Repo repository.VoiceChannelRepo
}

func (vc *VoiceChannels) Monitor(ctx context.Context) {
	updateTicker := time.NewTicker(vc.Interval)
	defer updateTicker.Stop()

	for {
		select {
		case <-updateTicker.C:
			vc.update()
		case <-ctx.Done():
			return
		}
	}
}

// update tracks how long 
func (vc *VoiceChannels) update() {
	ctx := context.Background()

	for _, guild := range vc.Session.State.Guilds {
		openSince, err := vc.Repo.GetVCOpenSince(ctx, guild.ID)
		if err != nil {
			fmt.Printf("Error: Failed to GetVCOpenDuration for guild %s. %v\n", guild.ID, err)
			continue
		}

		// Use to keep track if we have already recorded whether voice channel was active
		var voiceChannels map[string]bool = make(map[string]bool)
		for _, channel := range guild.Channels {
			if channel.Type == discordgo.ChannelTypeGuildVoice {
				voiceChannels[channel.ID] = false
			}
		}

		for _, voiceState := range guild.VoiceStates {
			if isDone, ok := voiceChannels[voiceState.ChannelID]; ok && isDone {
				continue
			} else {
				voiceChannels[voiceState.ChannelID] = true
			}

			timestamp, ok := openSince[voiceState.ChannelID]

			if ok {
				// If this voice chat has just been opened, set opening time
				if timestamp.IsZero() {
					openSince[voiceState.ChannelID] = time.Now()
				}
			} else {
				// Voice channel has not been tracked before
				openSince[voiceState.ChannelID] = time.Now()
			}
		}

		for vcId, isActive := range voiceChannels {
			if (!isActive) {
				// Set zero time value if voice channel is not active
				openSince[vcId] = time.Time{}
			}
		}

		err2 := vc.Repo.SetVCOpenSince(ctx, guild.ID, openSince)
		if err2 != nil {
			fmt.Printf("Error: Failed to SetVCOpenDuration for guild %s. %v\n", guild.ID, err2)
			continue
		}

		for channelID, timestamp := range openSince {
			channel, err := vc.Session.State.Channel(channelID)
			if err != nil {
				fmt.Printf("Error: Failed to edit channel %s for guild %s. %v\n", channelID, guild.ID, err)
				continue
			}

			// TODO: Find a better way to concat duration to channel name. What if the channel name has a ": " in it?
			newName := strings.Split(channel.Name, timeStampSeparator)[0]

			// If timestamp is non-zero (i.e. channel is active) then append timestamp
			if (!timestamp.IsZero()) {
				newName += timeStampSeparator + vc.formatDuration(time.Since(timestamp))
			} else {
				// If channel is inactive and timestamp has already been removed don't have to edit channel name again
				if newName == channel.Name {
					continue
				}
			}

			// Discord API applies rate-limiting to channel editing (2 requests per 10 minutes for each channel)
			_, err2 := vc.Session.ChannelEdit(channelID, newName)
			if err2 != nil {
				fmt.Printf("Error: Failed to edit channel %s for guild %s. %v\n", channelID, guild.ID, err2)
			}
		}
	}
}

func (vc *VoiceChannels) formatDuration(duration time.Duration) string {
	duration = duration.Round(time.Second)
	hour := duration / time.Hour
	duration -= hour * time.Hour
	minute := duration / time.Minute

	if hour == 0 {
		return fmt.Sprintf("%2dm", minute)
	} else {
		return fmt.Sprintf("%2dh%2dm", hour, minute)
	}
}