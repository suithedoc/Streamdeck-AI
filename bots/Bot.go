package bots

import (
	"OpenAITest/model"
	sd "OpenAITest/streamdeck"
	"OpenAITest/utils"
	"context"
	"fmt"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/micmonay/keybd_event"
	"github.com/muesli/streamdeck"
	"github.com/sashabaranov/go-openai"
	"golang.design/x/clipboard"
	"log"
	"time"
)

type AiBot struct {
	Name                               string
	SystemMsg                          string
	PromptMSg                          string
	StreamdeckDevice                   *streamdeck.Device
	StreamDeckButton                   int // Set to -1 to disable
	StreamDeckButtonWithHistory        int // Set to -1 to disable
	StreamDeckButtonWithHistoryAndCopy int // Set to -1 to disable
	streamdeckHandler                  *model.StreamdeckHandler
	OpenaiClient                       *openai.Client
	chatContent                        model.ChatContent
	speech                             *htgotts.Speech
	keyBonding                         *keybd_event.KeyBonding
}

func (bot *AiBot) init() {
	if bot.StreamDeckButton >= 0 {
		err := sd.SetStreamdeckButtonText(bot.StreamdeckDevice, uint8(bot.StreamDeckButton), bot.Name)
		if err != nil {
			log.Fatal(err)
		}
		bot.streamdeckHandler.AddOnPressHandler(bot.StreamDeckButton, func() error {
			go func() {
				isRecording = true
				utils.RecordAndSaveAudioAsMp3(fmt.Sprintf("audio%v.wav", bot.Name), quitChannel, finished)
			}()
			return nil
		})
		bot.streamdeckHandler.AddOnReleaseHandler(
			bot.StreamDeckButton,
			func() error {
				if isRecording {
					quitChannel <- true
					<-finished
					isRecording = false
					transcription, err := utils.ParseMp3ToText(fmt.Sprintf("audio%v.wav", bot.Name), bot.OpenaiClient)
					if err != nil {
						fmt.Printf("Error parsing mp3 to text 2: %s\n", err)
						return nil
					}
					EvaluateAssistantGptResponseStrings([]string{transcription}, false, bot.chatContent, bot.OpenaiClient, bot.speech)
				}
				return nil
			},
		)
	}
	if bot.StreamDeckButtonWithHistory >= 0 {
		err := sd.SetStreamdeckButtonText(bot.StreamdeckDevice, uint8(bot.StreamDeckButtonWithHistory), "HAssistant")
		if err != nil {
			log.Fatal(err)
		}
		bot.streamdeckHandler.AddOnPressHandler(int(bot.StreamDeckButtonWithHistory), func() error {
			go func() {
				isRecording = true
				utils.RecordAndSaveAudioAsMp3(fmt.Sprintf("audio%vHist.wav", bot.Name), quitChannel, finished)
			}()
			return nil
		})
		bot.streamdeckHandler.AddOnReleaseHandler(
			int(bot.StreamDeckButtonWithHistory),
			func() error {
				if isRecording {
					quitChannel <- true
					<-finished
					isRecording = false
					transcription, err := utils.ParseMp3ToText(fmt.Sprintf("audio%vHist.wav", bot.Name), bot.OpenaiClient)
					if err != nil {
						fmt.Printf("Error parsing mp3 to text 2: %s\n", err)
						return nil
					}
					EvaluateAssistantGptResponseStrings([]string{transcription}, true, bot.chatContent, bot.OpenaiClient, bot.speech)
				}
				return nil
			},
		)
	}

	if bot.StreamDeckButtonWithHistoryAndCopy >= 0 {
		err := sd.SetStreamdeckButtonText(bot.StreamdeckDevice, uint8(bot.StreamDeckButtonWithHistoryAndCopy), "HPAssistant")
		if err != nil {
			log.Fatal(err)
		}
		bot.streamdeckHandler.AddOnPressHandler(int(bot.StreamDeckButtonWithHistoryAndCopy), func() error {
			go func() {
				isRecording = true
				utils.RecordAndSaveAudioAsMp3(fmt.Sprintf("audio%vHistPaste.wav", bot.Name), quitChannel, finished)
			}()
			return nil
		})
		bot.streamdeckHandler.AddOnReleaseHandler(
			bot.StreamDeckButtonWithHistoryAndCopy,
			func() error {
				if !isRecording {
					return nil
				}
				quitChannel <- true
				<-finished
				isRecording = false
				transcription, err := utils.ParseMp3ToText(fmt.Sprintf("audio%vHistPaste.wav", bot.Name), bot.OpenaiClient)
				if err != nil {
					fmt.Printf("Error parsing mp3 to text 2: %s\n", err)
					return nil
				}
				respChan := clipboard.Watch(context.Background(), clipboard.FmtText)
				err = utils.CopySelectionToClipboard(bot.keyBonding)
				if err != nil {
					return err
				}

				// Wait for respChan 4 Seconds
				clipboardContent := ""
				var clipboardContentBytes []byte
				select {
				case <-time.After(4 * time.Second):
					log.Println("timeout waiting for clipboard")
					clipboardContentBytes = clipboard.Read(clipboard.FmtText)
				case clipboardContentBytes = <-respChan:
				}
				if clipboardContentBytes != nil {
					clipboardContent = string(clipboardContentBytes)
				}
				if string(clipboardContent) != "" {
					transcription = fmt.Sprintf("%s\n%s", transcription, string(clipboardContent))
				} else {
					err := bot.speech.Speak("Kein text im Clipboard gefunden")
					if err != nil {
						return err
					}
					return nil
				}

				EvaluateAssistantGptResponseStrings([]string{transcription}, true, bot.chatContent, bot.OpenaiClient, bot.speech)

				return nil
			})
	}
}

func (bot *AiBot) EvaluateGptResponseStringsWithHistory(requestLines []string) error {
	return EvaluateAssistantGptResponseStrings(requestLines, true, bot.chatContent, bot.OpenaiClient, bot.speech)
}

func (bot *AiBot) EvaluateGptResponseStrings(requestLines []string) error {
	return EvaluateAssistantGptResponseStrings(requestLines, false, bot.chatContent, bot.OpenaiClient, bot.speech)
}
