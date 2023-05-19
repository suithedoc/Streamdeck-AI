package bots

import (
	sd "OpenAITest/streamdeck"
	"OpenAITest/utils"
	"fmt"
	"log"
)

type AiBot struct {
	Name                        string
	SystemMsg                   string
	PromptMSg                   string
	StreamDeckButton            int
	StreamDeckButtonWithHistory int
}

func NewAiBot(name string, systemMsg string, promptMsg string, streamDeckButton int, streamDeckButtonWithHistory int) *AiBot {
	return &AiBot{
		Name:                        name,
		SystemMsg:                   systemMsg,
		PromptMSg:                   promptMsg,
		StreamDeckButton:            streamDeckButton,
		StreamDeckButtonWithHistory: streamDeckButtonWithHistory,
	}
}

func (aiBot *AiBot) Init() {
	if aiBot.StreamDeckButtonWithHistory >= 0 {
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
}
