package bot

import (
	"github.com/brandao07/vzcount/internal/data"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
)

func registerSlashCommands(discord *discordgo.Session) error {
	_, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "vzqueue",
		Description: "Start VampireZ queue monitoring",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "Channel to send alerts in",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "Role to ping",
				Required:    true,
			},
		},
	})

	if err != nil {
		return err
	}

	_, err = discord.ApplicationCommandCreate(discord.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "register",
		Description: "Register your Hypixel API key.",
	})

	if err != nil {
		return err
	}

	return nil
}

func monitorVampireZQueue(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	if i.ApplicationCommandData().Name == "vzqueue" {
		var channelID, roleID string

		userID := i.Member.User.ID

		for _, opt := range i.ApplicationCommandData().Options {
			switch opt.Name {
			case "channel":
				channelID = opt.ChannelValue(s).ID
			case "role":
				roleID = opt.RoleValue(s, i.GuildID).ID
			}
		}

		go VampireZQueue(s, userID, channelID, roleID, DefaultConfig())

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "🧛 Started VampireZ queue monitoring.",
			},
		})
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}
}

func register(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	if i.ApplicationCommandData().Name == "register" {
		modal := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "register_modal",
				Title:    "Register API Key",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "api_key_input",
								Label:       "Enter your Hypixel API key",
								Style:       discordgo.TextInputShort,
								Placeholder: "e.g. 01234567-89ab-cdef-0123-456789abcdef",
								Required:    true,
							},
						},
					},
				},
			},
		}
		err := s.InteractionRespond(i.Interaction, modal)
		if err != nil {
			log.Printf("Error sending modal: %v", err)
		}
	}
}

func handleModalSubmission(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionModalSubmit {
		return
	}

	if i.ModalSubmitData().CustomID == "register_modal" {
		userID := i.Member.User.ID
		var key string

		for _, c := range i.ModalSubmitData().Components {
			for _, inner := range c.(*discordgo.ActionsRow).Components {
				if input, ok := inner.(*discordgo.TextInput); ok && input.CustomID == "api_key_input" {
					key = input.Value
				}
			}
		}

		data.SetUser(userID, key, true)

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "✅ Your API key has been registered!",
				Flags:   discordgo.MessageFlagsEphemeral, // makes it visible only to the user
			},
		})
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}
}

func Run(token string) error {
	// create a session
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	// add a event handler
	discord.AddHandler(monitorVampireZQueue)
	discord.AddHandler(register)
	discord.AddHandler(handleModalSubmission)

	// open session
	err = discord.Open()
	if err != nil {
		return err
	}

	// close session
	defer func(discord *discordgo.Session) {
		err := discord.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(discord)

	// Register slash commands & handlers
	if err := registerSlashCommands(discord); err != nil {
		log.Fatalf("Failed to register slash commands: %v", err)
	}

	// keep bot running until there is NO os interruption (ctrl + C)
	log.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	return nil
}
