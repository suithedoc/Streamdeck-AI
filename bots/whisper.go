package bots

import (
	"OpenAITest/model"
	sd "OpenAITest/streamdeck"
	"OpenAITest/utils"
	"fmt"
	"github.com/micmonay/keybd_event"
	"github.com/muesli/streamdeck"
	"github.com/sashabaranov/go-openai"
	"golang.design/x/clipboard"
	"log"
	"time"
)

func TypeWhisperSTT(transcription string, kb *keybd_event.KeyBonding) error {
	clipboard.Write(clipboard.FmtText, []byte(transcription))
	time.Sleep(20 * time.Millisecond)
	kb.HasCTRL(true)
	kb.HasSHIFT(true)
	kb.SetKeys(keybd_event.VK_V)
	err := kb.Launching()
	if err != nil {
		fmt.Printf("Error launching keyboard: %v\n", err)
		return err
	}

	kb.HasCTRL(true)
	kb.HasSHIFT(false)
	kb.SetKeys(keybd_event.VK_V)
	err = kb.Launching()
	if err != nil {
		fmt.Printf("Error launching keyboard: %v\n", err)
		return err
	}

	return nil
}

func InitWhisperBot(streamdeckHandler *model.StreamdeckHandler, device *streamdeck.Device, kb *keybd_event.KeyBonding, client *openai.Client, button uint8) {
	err := sd.SetStreamdeckButtonText(device, button, "Whisper")
	if err != nil {
		log.Fatal(err)
	}
	streamdeckHandler.AddOnPressHandler(int(button), func() error {
		go func() {
			isRecording = true
			utils.RecordAndSaveAudioAsMp3("copyPaste.wav", quitChannel, finished)
		}()
		return nil
	})

	streamdeckHandler.AddOnReleaseHandler(int(button), func() error {
		if isRecording {
			quitChannel <- true
			<-finished
			isRecording = false
			transcription, err := utils.ParseMp3ToText("copyPaste.wav", client)
			if err != nil {
				fmt.Printf("Error parsing mp3 to text 2: %s\n", err)
				return nil
			}
			err = TypeWhisperSTT(transcription, kb)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
