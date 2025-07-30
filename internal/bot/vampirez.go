package bot

import (
	"fmt"
	"github.com/brandao07/vzcount/internal/data"
	"github.com/brandao07/vzcount/internal/hypixel"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

const (
	queueThreshold    = 9
	shortPollInterval = 5 * time.Second  // When VZ is not queueable
	longPollInterval  = 30 * time.Second // When VZ is queueable
)

type queueState int

const (
	stateNotQueueable queueState = iota
	stateQueueable
)

// Config holds configuration for VampireZ monitoring
type Config struct {
	QueueThreshold    int
	ShortPollInterval time.Duration
	LongPollInterval  time.Duration
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		QueueThreshold:    queueThreshold,
		ShortPollInterval: shortPollInterval,
		LongPollInterval:  longPollInterval,
	}
}

// VampireZQueue continuously monitors VampireZ player count and notifies on state changes
// Warning: Uses combined player count (in-game + lobby) from the Hypixel API
// Cannot distinguish between active players vs lobby queue, affecting accuracy
func VampireZQueue(discord *discordgo.Session, userID, channelID, roleID string, config Config) {
	var (
		currentState    = stateNotQueueable
		lastNotifyState = stateQueueable
		isFirstRun      = true
	)

	for {
		// I know shouldn't be here but like i'm too lazy
		user, ok := data.GetUser(userID)
		if !ok {
			_, _ = discord.ChannelMessageSend(channelID, "User is not registered, try the `/register` command.")
			return
		}
		client := &hypixel.API{
			Key: user.APIKey,
		}
		if !user.State {
			log.Printf("User %s is not enabled for VampireZ monitoring\n", userID)
			data.UpdateUserState(userID, true)
			return
		}
		count, err := client.VampireZCount()
		if err != nil {
			log.Printf("Error fetching VampireZ count: %v", err)
			time.Sleep(config.ShortPollInterval)
			continue
		}

		// Determine current queue state
		if count >= config.QueueThreshold {
			currentState = stateQueueable
		} else {
			currentState = stateNotQueueable
		}

		// Check if the state changed and notify accordingly
		if currentState != lastNotifyState || isFirstRun {
			var msg string
			switch currentState {
			case stateQueueable:
				log.Printf("VampireZ optimal to queue (%d players).\n", count)
				msg = fmt.Sprintf("<@&%s> 🧛 VampireZ is optimal to queue! (%d players online)", roleID, count)
			case stateNotQueueable:
				log.Printf("VampireZ no longer optimal to queue (%d players).\n", count)
				msg = fmt.Sprintf("🧛VampireZ is no longer optimal to queue. (%d players online)", count)
			}
			lastNotifyState = currentState
			isFirstRun = false
			// Send a message to Discord
			_, err := discord.ChannelMessageSend(channelID, msg)
			if err != nil {
				log.Printf("Error sending Discord message: %v", err)
			}
		}

		// Use the appropriate sleep interval based on the current state
		if currentState == stateQueueable {
			time.Sleep(config.LongPollInterval)
		} else {
			time.Sleep(config.ShortPollInterval)
		}
	}
}
