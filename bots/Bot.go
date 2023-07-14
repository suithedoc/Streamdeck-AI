package bots

import (
	"OpenAITest/model"
	sd "OpenAITest/streamdeck"
	"OpenAITest/utils"
	"context"
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/micmonay/keybd_event"
	"github.com/sashabaranov/go-openai"
	"golang.design/x/clipboard"
	"log"
	"strings"
	"time"
)

type AiBot struct {
	Name                   string
	SystemMsg              string
	PromptMSg              string
	StreamdeckDevice       sd.DeviceWrapper
	ButtonConfig           sd.StreamdeckButtonConfig
	streamdeckHandler      sd.IStreamdeckHandler
	OpenaiClient           *openai.Client
	ChatContent            model.ChatContent
	speech                 *htgotts.Speech
	keyBonding             *keybd_event.KeyBonding
	CompletionHistory      []openai.ChatCompletionMessage
	responseListeners      []func(*AiBot, string) error
	transcriptionListeners []func(*AiBot, string) error
	DisableAi              bool
}

func (bot *AiBot) AddResponseListener(listener func(*AiBot, string) error) {
	bot.responseListeners = append(bot.responseListeners, listener)
}

func (bot *AiBot) AddTranscriptionListener(listener func(*AiBot, string) error) {
	bot.transcriptionListeners = append(bot.transcriptionListeners, listener)
}

func (bot *AiBot) init() {
	if bot.ButtonConfig.ButtonIndex >= 0 {
		if bot.streamdeckHandler == nil {
			log.Fatal("Streamdeck handler not set")
		}
		err := bot.streamdeckHandler.AddButtonText(bot.ButtonConfig.Page, bot.ButtonConfig.ButtonIndex, bot.Name)
		if err != nil {
			log.Fatal(err)
		}
		bot.streamdeckHandler.AddOnPressHandler(bot.ButtonConfig.Page, bot.ButtonConfig.ButtonIndex, func() error {
			go func() {
				isRecording = true
				utils.RecordAndSaveAudioAsMp3(fmt.Sprintf("audio%v.wav", bot.Name), quitChannel, finished)
			}()
			return nil
		})
		bot.streamdeckHandler.AddOnReleaseHandler(bot.ButtonConfig.Page, bot.ButtonConfig.ButtonIndex, func() error {
			if isRecording {
				quitChannel <- true
				<-finished
				isRecording = false
				transcription, err := utils.ParseMp3ToText(fmt.Sprintf("audio%v.wav", bot.Name), bot.OpenaiClient)
				if err != nil {
					fmt.Printf("Error parsing mp3 to text 2: %s\n", err)
					return nil
				}
				for _, listener := range bot.transcriptionListeners {
					err = listener(bot, transcription)
					if err != nil {
						fmt.Printf("Error in transcription listener: %s\n", err)
					}
				}
				if bot.DisableAi == true {
					fmt.Printf("AI disabled for bot so ignore hostory gtp\n")
					return nil
				}
				err = bot.EvaluateGptResponseStrings([]string{transcription})
				if err != nil {
					fmt.Printf("Error evaluating gpt response: %s\n", err)
					return err
				}
			}
			return nil
		})
	}
	if bot.ButtonConfig.ButtonIndexHistory >= 0 {
		err := bot.streamdeckHandler.AddButtonText(bot.ButtonConfig.Page, bot.ButtonConfig.ButtonIndexHistory, "H"+bot.Name)
		if err != nil {
			log.Fatal(err)
		}
		bot.streamdeckHandler.AddOnPressHandler(bot.ButtonConfig.Page, bot.ButtonConfig.ButtonIndexHistory, func() error {
			go func() {
				isRecording = true
				utils.RecordAndSaveAudioAsMp3(fmt.Sprintf("audio%vHist.wav", bot.Name), quitChannel, finished)
			}()
			return nil
		})
		bot.streamdeckHandler.AddOnReleaseHandler(bot.ButtonConfig.Page, bot.ButtonConfig.ButtonIndexHistory, func() error {
			if isRecording {
				quitChannel <- true
				<-finished
				isRecording = false
				transcription, err := utils.ParseMp3ToText(fmt.Sprintf("audio%vHist.wav", bot.Name), bot.OpenaiClient)
				if err != nil {
					fmt.Printf("Error parsing mp3 to text 2: %s\n", err)
					return nil
				}
				for _, listener := range bot.transcriptionListeners {
					err = listener(bot, transcription)
					if err != nil {
						fmt.Printf("Error in transcription listener: %s\n", err)
					}
				}
				if bot.DisableAi == true {
					fmt.Printf("AI disabled for bot so ignore hostory gtp\n")
					return nil
				}
				err = bot.EvaluateGptResponseStringsWithHistory([]string{transcription})
				if err != nil {
					fmt.Printf("Error evaluating gpt response: %s\n", err)
					return err
				}
			}
			return nil
		})
	}

	if bot.ButtonConfig.ButtonIndexHistoryAndCopy >= 0 {
		err := bot.streamdeckHandler.AddButtonText(bot.ButtonConfig.Page, bot.ButtonConfig.ButtonIndexHistoryAndCopy, "HP"+bot.Name)
		if err != nil {
			log.Fatal(err)
		}
		bot.streamdeckHandler.AddOnPressHandler(bot.ButtonConfig.Page, bot.ButtonConfig.ButtonIndexHistoryAndCopy, func() error {
			go func() {
				isRecording = true
				utils.RecordAndSaveAudioAsMp3(fmt.Sprintf("audio%vHistPaste.wav", bot.Name), quitChannel, finished)
			}()
			return nil
		})
		bot.streamdeckHandler.AddOnReleaseHandler(bot.ButtonConfig.Page, bot.ButtonConfig.ButtonIndexHistoryAndCopy, func() error {
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
			for _, listener := range bot.transcriptionListeners {
				err = listener(bot, transcription)
				if err != nil {
					fmt.Printf("Error in transcription listener: %s\n", err)
				}
			}
			if bot.DisableAi == true {
				fmt.Printf("AI disabled for bot so ignore hostory gtp\n")
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
			err = bot.EvaluateGptResponseStringsWithHistory([]string{transcription})
			if err != nil {
				fmt.Printf("Error evaluating gpt response: %s\n", err)
				return err
			}
			return nil
		})
	}
}

func (bot *AiBot) EvaluateAndReturnGptResponse(input []string, withHistory bool, chatContent model.ChatContent, client *openai.Client) (string, error) {
	joinedRequestMessage := strings.Join(input, "\n")
	if withHistory {
		chatContent.HistoryMessages = bot.CompletionHistory
	} else {
		bot.CompletionHistory = []openai.ChatCompletionMessage{}
	}
	answer, err := utils.SendChatRequest(joinedRequestMessage, chatContent, client)
	if err != nil {
		log.Fatal(err)
	}
	bot.CompletionHistory = append(bot.CompletionHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: joinedRequestMessage,
	})
	bot.CompletionHistory = append(bot.CompletionHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: answer,
	})
	result := markdown.Render(answer, 80, 6)
	fmt.Println(string(result))

	return answer, nil
}

func (bot *AiBot) EvaluateGptResponseStringsWithHistory(requestLines []string) error {
	response, err := bot.EvaluateAndReturnGptResponse(requestLines, true, bot.ChatContent, bot.OpenaiClient)
	if err != nil {
		return err
	}
	for _, responseListener := range bot.responseListeners {
		responseListener(bot, response)
	}
	return nil
}

func (bot *AiBot) EvaluateGptResponseStrings(requestLines []string) error {
	response, err := bot.EvaluateAndReturnGptResponse(requestLines, false, bot.ChatContent, bot.OpenaiClient)
	if err != nil {
		return err
	}
	for _, responseListener := range bot.responseListeners {
		responseListener(bot, response)
	}
	return nil
}

func (bot *AiBot) EvaluateAndReturn(requestLines []string) (string, error) {
	return bot.EvaluateAndReturnGptResponse(requestLines, false, bot.ChatContent, bot.OpenaiClient)
}

func (bot *AiBot) EvaluateAndReturnWithHistory(requestLines []string) (string, error) {
	return bot.EvaluateAndReturnGptResponse(requestLines, true, bot.ChatContent, bot.OpenaiClient)
}
