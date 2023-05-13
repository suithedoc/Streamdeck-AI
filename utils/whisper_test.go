package utils

import (
	"github.com/tinyzimmer/go-gst/gst"
	"testing"
)

func TestRecordAndSaveAudioAsWav(t *testing.T) {
	gst.Init(nil)

	finishedChan := make(chan bool)
	RecordAndSaveAudioAsWav("test.wav", finishedChan)
	<-finishedChan
}
