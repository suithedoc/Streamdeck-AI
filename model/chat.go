package model

import (
	"github.com/sashabaranov/go-openai"
)

type ChatContent struct {
	SystemMsg       string
	PromptMsg       string
	HistoryMessages []openai.ChatCompletionMessage
}
