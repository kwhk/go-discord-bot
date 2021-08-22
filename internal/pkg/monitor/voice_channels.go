package monitor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kwhk/go-discord-bot/internal/pkg/repository"
)

const (
	timeStampSeparator = " - "
)

type VoiceChannels struct {
	Session           *discordgo.Session
	Interval          time.Duration
	VoiceChannelRepo  repository.VoiceChannelRepo
	NotificationsRepo repository.NotificationsRepo
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
		err := vc.monitorVoiceChannelDuration(ctx, guild)
		if err != nil {
			continue
		}
	}
}

func (vc *VoiceChannels) monitorVoiceChannelDuration(ctx context.Context, guild *discordgo.Guild) error {
	openSince, err := vc.VoiceChannelRepo.GetVCOpenSince(ctx, guild.ID)
	if err != nil {
		fmt.Printf("Error: Failed to GetVCOpenDuration for guild %s. %v\n", guild.ID, err)
		return err
	}

	// Use to keep track if we have already recorded whether voice channel was active
	var voiceChannels map[string]bool = make(map[string]bool)
	for _, channel := range guild.Channels {
		if channel.Type == discordgo.ChannelTypeGuildVoice {
			voiceChannels[channel.ID] = false
		}
	}

	for _, voiceState := range guild.VoiceStates {
		// If already checked voice channel then skip
		if isDone, ok := voiceChannels[voiceState.ChannelID]; ok && isDone {
			continue
		}

		// Set voice channel to be active
		voiceChannels[voiceState.ChannelID] = true
		timestamp, ok := openSince[voiceState.ChannelID]

		// If voice chat has been just been opened OR first time voice channel
		// is being tracked then set timestamp
		if (ok && timestamp.IsZero()) || !ok {
			openSince[voiceState.ChannelID] = time.Now()
		}
	}

	for vcId, isActive := range voiceChannels {
		if !isActive {
			// If voice channel has just been closed, notify channels
			if timestamp, ok := openSince[vcId]; ok && !timestamp.IsZero() {
				duration := time.Since(timestamp)
				channel, err := vc.Session.State.Channel(vcId)
				if err != nil {
					fmt.Printf("Error: Failed to get channel %s\n", vcId)
					continue
				}

				channelName := strings.Split(channel.Name, timeStampSeparator)[0]
				if duration > 1*time.Second {
					listeners, err := vc.NotificationsRepo.GetNotificationChannels(ctx, guild.ID)
					if err == nil {
						for _, listener := range listeners {
							msg := fmt.Sprintf("Call ended in channel \"%s\" after %s", channelName, vc.formatDuration(duration))
							vc.Session.ChannelMessageSend(listener, msg)
						}
					}
				}
			}

			// Set zero time value if voice channel is not active
			openSince[vcId] = time.Time{}
		}
	}

	err2 := vc.VoiceChannelRepo.SetVCOpenSince(ctx, guild.ID, openSince)
	if err2 != nil {
		fmt.Printf("Error: Failed to SetVCOpenDuration for guild %s. %v\n", guild.ID, err2)
		return err2
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
		if !timestamp.IsZero() {
			newName += timeStampSeparator + vc.formatDuration(time.Since(timestamp))
		} else {
			// If channel is inactive and timestamp has already been removed don't have to edit channel name again
			if newName == channel.Name {
				continue
			}
		}

		// Discord API applies rate-limiting to channel editing (2 requests per 10 minutes for each channel)
		_, err2 := vc.Session.ChannelEdit(channelID, newName)
		fmt.Println(err2)
		if err2 != nil {
			fmt.Printf("Error: Failed to edit channel %s for guild %s. %v\n", channelID, guild.ID, err2)
		}
	}

	return nil
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
