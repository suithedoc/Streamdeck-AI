package model

import "github.com/sashabaranov/go-openai"

type ChatContent struct {
	SystemMsg       string
	PromptMsg       string
	HistoryMessages []openai.ChatCompletionMessage
}

type StreamdeckHandler struct {
	streamDeckButtonIdToOnPressHandler   map[int]func() error
	streamDeckButtonIdToOnReleaseHandler map[int]func() error
}

func NewStreamdeckHandler() *StreamdeckHandler {
	return &StreamdeckHandler{
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
