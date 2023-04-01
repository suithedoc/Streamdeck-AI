package bots

import (
	"OpenAITest/model"
	"OpenAITest/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sashabaranov/go-openai"
	"log"
	"time"
)

func StartDiscordBot(chatContent model.ChatContent, properties map[string]string, client *openai.Client) error {
	fmt.Printf("starting Discord-Bot\n")
	discordBot, err := discordgo.New("Bot " + properties["discordBotToken"])
	if err != nil {
		fmt.Printf("Cant create a new DiscordBot-Session: %v", err.Error())
		return nil
	}
	discordBot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		channel, _ := s.Channel(m.ChannelID)
		fmt.Printf("\nChannel: %v \nAuthor: %v\nMessage: %v \n", channel.Name, m.Author, m.Content)
		if m.Content == "" {
			fmt.Println("Message content is empty")
		}
		answer, err := utils.SendChatRequest(m.Content, chatContent, client)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Millisecond * 500)
		_, err = s.ChannelMessageSend(m.ChannelID, answer)
		if err != nil {
			return
		}
		fmt.Println("Answer: ", answer)

	})

	discordBot.Identify.Intents = discordgo.IntentsGuildMessages
	err = discordBot.Open()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Discord-Bot is online")
	return nil
}

func InitDiscordBot(client *openai.Client, properties map[string]string) *model.ChatContent {
	discordSystemMsg := properties["discordSystemMsg"]
	discordPromptMsg := properties["discordPromptMsg"]
	discordChatContent := model.ChatContent{
		SystemMsg:       discordSystemMsg,
		PromptMsg:       discordPromptMsg,
		HistoryMessages: nil,
	}

	go func() {
		err := StartDiscordBot(discordChatContent, properties, client)
		if err != nil {
			fmt.Printf("Stopped DiscordBot: %v", err.Error())
		}
	}()
	return nil
}
