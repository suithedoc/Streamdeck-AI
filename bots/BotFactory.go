package bots

import (
	"OpenAITest/model"
	sd "OpenAITest/streamdeck"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/micmonay/keybd_event"
	"github.com/sashabaranov/go-openai"
)

type BotFactory struct {
	streamdeckHandler sd.IStreamdeckHandler
	OpenaiClient      *openai.Client
	StreamdeckDevice  sd.DeviceWrapper
	speeches          map[string]*htgotts.Speech
	keyBonding        *keybd_event.KeyBonding
}

func NewBotFactory(streamdeckHandler sd.IStreamdeckHandler,
	openaiClient *openai.Client, streamdeckDevice sd.DeviceWrapper,
	keyBonding *keybd_event.KeyBonding) *BotFactory {
	return &BotFactory{
		streamdeckHandler: streamdeckHandler,
		OpenaiClient:      openaiClient,
		StreamdeckDevice:  streamdeckDevice,
		speeches:          make(map[string]*htgotts.Speech),
		keyBonding:        keyBonding,
	}
}

func (bf *BotFactory) createBaseBot(name string, systemMsg string, promptMsg string,
	streamdeckButtonConfig sd.StreamdeckButtonConfig, speechLanguage string) *AiBot {
	aiBot := &AiBot{
		Name:              name,
		SystemMsg:         systemMsg,
		PromptMSg:         promptMsg,
		ButtonConfig:      streamdeckButtonConfig,
		StreamdeckDevice:  bf.StreamdeckDevice,
		streamdeckHandler: bf.streamdeckHandler,
		OpenaiClient:      bf.OpenaiClient,
		ChatContent: model.ChatContent{
			SystemMsg:       systemMsg,
			PromptMsg:       promptMsg,
			HistoryMessages: make([]openai.ChatCompletionMessage, 0),
		},
		keyBonding:             bf.keyBonding,
		CompletionHistory:      make([]openai.ChatCompletionMessage, 0),
		responseListeners:      make([]func(*AiBot, string) error, 0),
		transcriptionListeners: make([]func(*AiBot, string) error, 0),
		DisableAi:              false,
	}
	speech, ok := bf.speeches[speechLanguage]
	if !ok {
		speech = &htgotts.Speech{Folder: "audio", Language: speechLanguage, Handler: &handlers.Native{}}
		bf.speeches[speechLanguage] = speech
	}
	aiBot.speech = speech
	return aiBot
}

//func (bf *BotFactory) CreateBot(name string, systemMsg string, promptMsg string,
//	streamDeckButton int, speechLanguage string) *AiBot {
//	aiBot := bf.createBaseBot(name, systemMsg, promptMsg, streamDeckButton, speechLanguage)
//	aiBot.init()
//	return aiBot
//}

func (bf *BotFactory) CreateBot(name string, systemMsg string, promptMsg string,
	buttonConfig sd.StreamdeckButtonConfig, speechLanguage string) *AiBot {
	aiBot := bf.createBaseBot(name, systemMsg, promptMsg, buttonConfig, speechLanguage)
	aiBot.init()
	return aiBot
}

//func (bf *BotFactory) CreateBotWithHistory(name string, systemMsg string, promptMsg string,
//	streamDeckButton int, streamDeckButtonWithHistory int, speechLanguage string) *AiBot {
//	aiBot := bf.createBaseBot(name, systemMsg, promptMsg, streamDeckButton, speechLanguage)
//	aiBot.StreamDeckButtonWithHistory = streamDeckButtonWithHistory
//	aiBot.init()
//	return aiBot
//}
//
//func (bf *BotFactory) CreateBotWithCopy(name string, systemMsg string, promptMsg string,
//	streamDeckButton int, streamDeckButtonWithHistoryAndCopy int, speechLanguage string) *AiBot {
//	aiBot := bf.createBaseBot(name, systemMsg, promptMsg, streamDeckButton, speechLanguage)
//	aiBot.StreamDeckButtonWithHistoryAndCopy = streamDeckButtonWithHistoryAndCopy
//	aiBot.init()
//	return aiBot
//}
//
//func (bf *BotFactory) CreateBotWithHistoryAndCopy(name string, systemMsg string, promptMsg string,
//	streamDeckButton int, streamDeckButtonWithHistory int, streamDeckButtonWithHistoryAndCopy int, speechLanguage string) *AiBot {
//	aiBot := bf.createBaseBot(name, systemMsg, promptMsg, streamDeckButton, speechLanguage)
//	aiBot.StreamDeckButtonWithHistory = streamDeckButtonWithHistory
//	aiBot.StreamDeckButtonWithHistoryAndCopy = streamDeckButtonWithHistoryAndCopy
//	aiBot.init()
//	return aiBot
//}
