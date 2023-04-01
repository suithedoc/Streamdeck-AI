package utils

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"github.com/tinyzimmer/go-gst/gst"
	"log"
	"time"
)

func ParseMp3ToText(mp3FileName string, client *openai.Client) (string, error) {
	audioRequest := openai.AudioRequest{
		Model:    "whisper-1",
		FilePath: mp3FileName,
	}
	transcription, err := client.CreateTranscription(context.Background(), audioRequest)
	if err != nil {
		return "", err
	}
	log.Printf("Transcription: %s\n", transcription)
	return transcription.Text, nil
}

func RecordAndSaveAudioAsMp3(mp3FileName string, quit chan bool, finished chan bool) {
	pielineString := "autoaudiosrc ! audioconvert ! audioresample ! audio/x-raw,rate=16000,channels=1,format=S16LE ! wavenc ! filesink location=" + mp3FileName
	pipeline, err := gst.NewPipelineFromString(pielineString)
	if err != nil {
		log.Fatal(err)
	}
	err = pipeline.SetState(gst.StatePlaying)
	if err != nil {
		log.Fatal(err)
	}
	for {
		switch {
		case <-quit:
			fmt.Println("Sending EOS")
			pipeline.SendEvent(gst.NewEOSEvent())

			msg := pipeline.GetBus().TimedPop(3 * time.Second)
			switch msg.Type() {
			case gst.MessageEOS:
				fmt.Println("Got EOS")
			case gst.MessageError:
				fmt.Println("Got Error")
			default:
			}
			if msg == nil {
				log.Println("Timed out waiting for EOS message")
			} else {
				log.Printf("Received %s message\n", msg.Type())
			}

			err := pipeline.SetState(gst.StateNull)
			if err != nil {
				log.Fatal(err)
			}
			finished <- true
			fmt.Println("Sleeping")
			return
		}
	} //use googler --noprompt to search for information about the current weather and then try to download and display the content of the first result in the terminal
}
