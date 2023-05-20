package model

import (
	"github.com/muesli/streamdeck"
	"testing"
)

func TestStreamdeckHandler_TraverseCoslAndRowByButtonId(t *testing.T) {
	device, err := streamdeck.Devices()
	if err != nil {
		t.Fatal(err)
	}
	if len(device) == 0 {
		t.Fatal("No streamdeck found")
	}
	streamdeckHandler := NewStreamdeckHandler(&device[0])
	buttonId := 0
	newButtonId := streamdeckHandler.TraverseButtonId(buttonId)
	if newButtonId != 0 {
		t.Fatalf("Expected buttonId to be 0, but got %v", newButtonId)
	}

	buttonId = 1
	newButtonId = streamdeckHandler.TraverseButtonId(buttonId)
	if newButtonId != 5 {
		t.Fatalf("Expected buttonId to be 5, but got %v", newButtonId)
	}

	buttonId = 2
	newButtonId = streamdeckHandler.TraverseButtonId(buttonId)
	if newButtonId != 10 {
		t.Fatalf("Expected buttonId to be 10, but got %v", newButtonId)
	}

	buttonId = 3
	newButtonId = streamdeckHandler.TraverseButtonId(buttonId)
	if newButtonId != 1 {
		t.Fatalf("Expected buttonId to be 1, but got %v", newButtonId)
	}

	buttonId = 4
	newButtonId = streamdeckHandler.TraverseButtonId(buttonId)
	if newButtonId != 6 {
		t.Fatalf("Expected buttonId to be 6, but got %v", newButtonId)
	}

	buttonId = 5
	newButtonId = streamdeckHandler.TraverseButtonId(buttonId)
	if newButtonId != 11 {
		t.Fatalf("Expected buttonId to be 11, but got %v", newButtonId)
	}
}
