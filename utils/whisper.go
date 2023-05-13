package utils

import "C"
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

//type LevelMessage struct {
//	Name           string
//	RunningTime    int64
//	Rms            int `gst:"rms"`
//	RMSDecibel     float64
//	Peak           float64
//	PeakDecibel    float64
//	Decay          float64
//	DecayDecibel   float64
//	LastMessage    bool
//	Overrun        bool
//	CaptureStopped bool
//}

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
	}
}

func RecordAndSaveAudioAsWav(mp3FileName string, timeout time.Duration, finished chan bool) {
	pielineString := "autoaudiosrc ! audioconvert ! audioresample ! removesilence remove=true silent=false ! audio/x-raw,rate=16000,channels=1,format=S16LE ! wavenc ! filesink location=" + mp3FileName
	pipeline, err := gst.NewPipelineFromString(pielineString)
	if err != nil {
		log.Fatal(err)
	}
	err = pipeline.SetState(gst.StatePlaying)
	if err != nil {
		log.Fatal(err)
	}

	bus := pipeline.GetBus()
	var silenceDetectedTimestamp, silenceFinishedTimestamp uint64

	for {
		msg := bus.TimedPop(3 * time.Second)
		if msg == nil {
			continue
		}

		switch msg.Type() {
		case gst.MessageError:
			fmt.Println("Got Error")
			err := pipeline.SetState(gst.StateNull)
			if err != nil {
				log.Fatal(err)
			}
			finished <- false
			return
		case gst.MessageElement:
			st := msg.GetStructure()
			if st == nil {
				continue
			}

			if st.Name() == "removesilence" {
				if ts, exists := st.Values()["silence_detected"]; exists {
					silenceDetectedTimestamp = ts.(uint64)
					//fmt.Printf("Silence detected timestamp: %d\n", silenceDetectedTimestamp)
				} else if ts, exists := st.Values()["silence_finished"]; exists {
					silenceFinishedTimestamp = ts.(uint64)
					//fmt.Printf("Silence finished timestamp: %d\n", silenceFinishedTimestamp)
				}

				if silenceFinishedTimestamp > silenceDetectedTimestamp && (silenceFinishedTimestamp-silenceDetectedTimestamp) >= uint64(timeout) {
					pipeline.SendEvent(gst.NewEOSEvent())
				}
			}
		case gst.MessageEOS:
			fmt.Println("Got EOS")
			fmt.Println("Setting State to Null")
			err := pipeline.SetState(gst.StateNull)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Did set state to null")
			finished <- true
			fmt.Println("Sleeping")
			return
		default:
		}
	}
}

//func RecordAndSaveAudioAsMp3(mp3FileName string, quit chan bool, finished chan bool) {
//	pielineString := "autoaudiosrc ! audioconvert ! audioresample ! audio/x-raw,rate=16000,channels=1,format=S16LE ! wavenc ! filesink location=" + mp3FileName
//	pipeline, err := gst.NewPipelineFromString(pielineString)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = pipeline.SetState(gst.StatePlaying)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for {
//		switch {
//		case <-quit:
//			fmt.Println("Sending EOS")
//			pipeline.SendEvent(gst.NewEOSEvent())
//
//			msg := pipeline.GetBus().TimedPop(3 * time.Second)
//			switch msg.Type() {
//			case gst.MessageEOS:
//				fmt.Println("Got EOS")
//			case gst.MessageError:
//				fmt.Println("Got Error")
//			default:
//			}
//			if msg == nil {
//				log.Println("Timed out waiting for EOS message")
//			} else {
//				log.Printf("Received %s message\n", msg.Type())
//			}
//
//			err := pipeline.SetState(gst.StateNull)
//			if err != nil {
//				log.Fatal(err)
//			}
//			finished <- true
//			fmt.Println("Sleeping")
//			return
//		}
//	}
//}
