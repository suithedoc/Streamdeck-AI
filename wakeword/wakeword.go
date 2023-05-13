package wakeword

import (
	porcupine "github.com/Picovoice/porcupine/binding/go/v2"
	pvrecorder "github.com/Picovoice/pvrecorder/sdk/go"
	"log"
)

func StartListeningToWakeword(keyword string, sensitivity float32, accessKey string, wordDetectedChannel chan bool) error {

	p := porcupine.Porcupine{}
	p.AccessKey = accessKey
	p.Sensitivities = []float32{sensitivity}

	buildInKeyword := porcupine.BuiltInKeyword(keyword)
	p.BuiltInKeywords = []porcupine.BuiltInKeyword{buildInKeyword}
	err := p.Init()
	if err != nil {
		return err
	}
	defer func() {
		err := p.Delete()
		if err != nil {
			log.Fatalf("failed to delete porcupine: %+v", err)
		}
	}()

	recorder := pvrecorder.PvRecorder{
		DeviceIndex:    -1,
		FrameLength:    porcupine.FrameLength,
		BufferSizeMSec: 1000,
		LogOverflow:    0,
	}

	err = recorder.Init()
	if err != nil {
		return err
	}
	defer recorder.Delete()

	log.Printf("Using device: %s", recorder.GetSelectedDevice())
	err = recorder.Start()
	if err != nil {
		return err
	}
	log.Printf("listening...")

	for {
		audioFrame, err := recorder.Read()
		if err != nil {
			return err
		}
		wordIndex, err := p.Process(audioFrame)
		if err != nil {
			log.Printf("error: %+v", err)
			continue
		}
		if wordIndex >= 0 {
			log.Printf("detected word")
			wordDetectedChannel <- true
		}
	}

}
