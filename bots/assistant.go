package bots

import (
	"OpenAITest/model"
	sd "OpenAITest/streamdeck"
	"OpenAITest/utils"
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/muesli/streamdeck"
	"github.com/sashabaranov/go-openai"
	"log"
	"strings"
)

var assistantCompletionHistory []openai.ChatCompletionMessage

func EvaluateAssistantGptResponseStrings(input []string, withHistory bool, chatContent model.ChatContent, client *openai.Client, speech *htgotts.Speech) error {
	joinedRequestMessage := strings.Join(input, "\n")
	if withHistory {
		chatContent.HistoryMessages = assistantCompletionHistory
	} else {
		assistantCompletionHistory = []openai.ChatCompletionMessage{}
	}
	answer, err := utils.SendChatRequest(joinedRequestMessage, chatContent, client)
	if err != nil {
		log.Fatal(err)
	}
	assistantCompletionHistory = append(assistantCompletionHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: joinedRequestMessage,
	})
	assistantCompletionHistory = append(assistantCompletionHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: answer,
	})
	result := markdown.Render(answer, 80, 6)
	fmt.Println(string(result))

	chunks := utils.SplitStringIntoWordChunks(answer, 100)
	for _, chunk := range chunks {
		err = speech.Speak(chunk)
		if err != nil {
			return err
		}
	}
	return nil
}

func InitAssistantGPTBot(client *openai.Client, device *streamdeck.Device, properties map[string]string,
	streamdeckHandler *model.StreamdeckHandler, speech *htgotts.Speech, buttonWithoutHistory int16, buttonWithHistory int16) *model.ChatContent {
	assistantCompletionHistory = []openai.ChatCompletionMessage{}
	assistantChatContent := model.ChatContent{
		SystemMsg:       properties["assistantSystemMsg"],
		PromptMsg:       properties["assistantPromptMsg"],
		HistoryMessages: []openai.ChatCompletionMessage{},
	}

	if buttonWithoutHistory >= 0 {
		err := sd.SetStreamdeckButtonText(device, uint8(buttonWithoutHistory), "Assistant")
		if err != nil {
			log.Fatal(err)
		}
		streamdeckHandler.AddOnPressHandler(int(buttonWithoutHistory), func() error {
			go func() {
				isRecording = true
				utils.RecordAndSaveAudioAsMp3("audioAssist.wav", quitChannel, finished)
			}()
			return nil
		})
		streamdeckHandler.AddOnReleaseHandler(
			int(buttonWithoutHistory),
			func() error {
				if isRecording {
					quitChannel <- true
					<-finished
					isRecording = false
					transcription, err := utils.ParseMp3ToText("audioAssist.wav", client)
					if err != nil {
						fmt.Printf("Error parsing mp3 to text 2: %s\n", err)
						return nil
					}
					EvaluateAssistantGptResponseStrings([]string{transcription}, false, assistantChatContent, client, speech)
				}
				return nil
			},
		)
	}

	if buttonWithHistory >= 0 {
		err := sd.SetStreamdeckButtonText(device, uint8(buttonWithHistory), "HAssistant")
		if err != nil {
			log.Fatal(err)
		}
		streamdeckHandler.AddOnPressHandler(int(buttonWithHistory), func() error {
			go func() {
				isRecording = true
				utils.RecordAndSaveAudioAsMp3("audioAssistHist.wav", quitChannel, finished)
			}()
			return nil
		})
		streamdeckHandler.AddOnReleaseHandler(
			int(buttonWithHistory),
			func() error {
				if isRecording {
					quitChannel <- true
					<-finished
					isRecording = false
					transcription, err := utils.ParseMp3ToText("audioAssistHist.wav", client)
					if err != nil {
						fmt.Printf("Error parsing mp3 to text 2: %s\n", err)
						return nil
					}
					EvaluateAssistantGptResponseStrings([]string{transcription}, true, assistantChatContent, client, speech)
				}
				return nil
			},
		)
	}

	return &assistantChatContent
}
