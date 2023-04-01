package utils

import (
	"OpenAITest/model"
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
)

var (
	OpenaiModel string
)

func SendChatRequest(text string, chatContent model.ChatContent, client *openai.Client) (string, error) {

	fmt.Printf("Sending chatContent. text:%v \n", text)
	fmt.Println("chatHistory:")
	for _, msg := range chatContent.HistoryMessages {
		fmt.Printf("role: %v, content: %v \n", msg.Role, msg.Content)
	}

	completionMessages := make([]openai.ChatCompletionMessage, 0)
	completionMessages = append(completionMessages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: chatContent.SystemMsg,
	})
	if chatContent.HistoryMessages != nil && len(chatContent.HistoryMessages) > 0 {
		completionMessages = append(completionMessages, chatContent.HistoryMessages...)
	}
	completionMessages = append(completionMessages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: chatContent.PromptMsg + text,
	})

	modelToUse := openai.GPT3Dot5Turbo
	switch OpenaiModel {
	case "gpt3.5":
		modelToUse = openai.GPT3Dot5Turbo
	case "gpt4":
		modelToUse = openai.GPT4
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    modelToUse,
			Messages: completionMessages,
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
