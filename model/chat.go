package model

import (
	"fmt"
	"github.com/muesli/streamdeck"
	"github.com/sashabaranov/go-openai"
	"log"
)

type ChatContent struct {
	SystemMsg       string
	PromptMsg       string
	HistoryMessages []openai.ChatCompletionMessage
}

type StreamdeckHandler struct {
	device                               *streamdeck.Device
	streamDeckButtonIdToOnPressHandler   map[int]func() error
	streamDeckButtonIdToOnReleaseHandler map[int]func() error
}

func NewStreamdeckHandler(device *streamdeck.Device) *StreamdeckHandler {
	return &StreamdeckHandler{
		device:                               device,
		streamDeckButtonIdToOnPressHandler:   make(map[int]func() error),
		streamDeckButtonIdToOnReleaseHandler: make(map[int]func() error),
	}
}

func (sh *StreamdeckHandler) AddOnPressHandler(buttonId int, handler func() error) {
	sh.streamDeckButtonIdToOnPressHandler[buttonId] = handler
}

func (sh *StreamdeckHandler) AddOnReleaseHandler(buttonId int, handler func() error) {
	sh.streamDeckButtonIdToOnReleaseHandler[buttonId] = handler
}

func (sh *StreamdeckHandler) GetOnPressHandler(buttonId int) (func() error, bool) {
	handler, ok := sh.streamDeckButtonIdToOnPressHandler[buttonId]
	return handler, ok
}

func (sh *StreamdeckHandler) GetOnReleaseHandler(buttonId int) (func() error, bool) {
	handler, ok := sh.streamDeckButtonIdToOnReleaseHandler[buttonId]
	return handler, ok
}

func (sh *StreamdeckHandler) StartAsync() {
	go func() {
		keys, err := sh.device.ReadKeys()
		if err != nil {
			log.Fatal(err)
		}
		for {
			select {
			case key := <-keys:
				fmt.Printf("Key pressed index %v, is pressed %v\n", key.Index, key.Pressed)
				if key.Pressed {
					if handler, ok := sh.GetOnPressHandler(int(key.Index)); ok {
						err := handler()
						if err != nil {
							log.Fatal(err)
						}
					}
				} else {
					if handler, ok := sh.GetOnReleaseHandler(int(key.Index)); ok {
						err := handler()
						if err != nil {
							log.Fatal(err)
						}
					}
				}
			}
		}
	}()
}
