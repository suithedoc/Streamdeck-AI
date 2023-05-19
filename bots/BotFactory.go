package bots

import (
	"OpenAITest/model"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/micmonay/keybd_event"
	"github.com/muesli/streamdeck"
	"github.com/sashabaranov/go-openai"
)

type BotFactory struct {
	streamdeckHandler *model.StreamdeckHandler
	OpenaiClient      *openai.Client
	StreamdeckDevice  *streamdeck.Device
	speeches          map[string]*htgotts.Speech
	keyBonding        *keybd_event.KeyBonding
}

func NewBotFactory(streamdeckHandler *model.StreamdeckHandler,
	openaiClient *openai.Client, streamdeckDevice *streamdeck.Device,
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
	streamDeckButton int, speechLanguage string) *AiBot {
	aiBot := &AiBot{
		Name:                               name,
		SystemMsg:                          systemMsg,
		PromptMSg:                          promptMsg,
		StreamDeckButton:                   streamDeckButton,
		StreamDeckButtonWithHistory:        -1,
		StreamDeckButtonWithHistoryAndCopy: -1,
		StreamdeckDevice:                   bf.StreamdeckDevice,
		streamdeckHandler:                  bf.streamdeckHandler,
		OpenaiClient:                       bf.OpenaiClient,
		ChatContent: model.ChatContent{
			SystemMsg:       systemMsg,
			PromptMsg:       promptMsg,
			HistoryMessages: make([]openai.ChatCompletionMessage, 0),
		},
		keyBonding:        bf.keyBonding,
		CompletionHistory: make([]openai.ChatCompletionMessage, 0),
	}
	speech, ok := bf.speeches[speechLanguage]
	if !ok {
		speech = &htgotts.Speech{Folder: "audio", Language: speechLanguage, Handler: &handlers.Native{}}
		bf.speeches[speechLanguage] = speech
	}
	aiBot.speech = speech
	return aiBot
}

func (bf *BotFactory) CreateBot(name string, systemMsg string, promptMsg string,
	streamDeckButton int, speechLanguage string) *AiBot {
	aiBot := bf.createBaseBot(name, systemMsg, promptMsg, streamDeckButton, speechLanguage)
	aiBot.init()
	return aiBot
}

func (bf *BotFactory) CreateBotWithHistory(name string, systemMsg string, promptMsg string,
	streamDeckButton int, streamDeckButtonWithHistory int, speechLanguage string) *AiBot {
	aiBot := bf.createBaseBot(name, systemMsg, promptMsg, streamDeckButton, speechLanguage)
	aiBot.StreamDeckButtonWithHistory = streamDeckButtonWithHistory
	aiBot.init()
	return aiBot
}

func (bf *BotFactory) CreateBotWithCopy(name string, systemMsg string, promptMsg string,
	streamDeckButton int, streamDeckButtonWithHistoryAndCopy int, speechLanguage string) *AiBot {
	aiBot := bf.createBaseBot(name, systemMsg, promptMsg, streamDeckButton, speechLanguage)
	aiBot.StreamDeckButtonWithHistoryAndCopy = streamDeckButtonWithHistoryAndCopy
	aiBot.init()
	return aiBot
}

func (bf *BotFactory) CreateBotWithHistoryAndCopy(name string, systemMsg string, promptMsg string,
	streamDeckButton int, streamDeckButtonWithHistory int, streamDeckButtonWithHistoryAndCopy int, speechLanguage string) *AiBot {
	aiBot := bf.createBaseBot(name, systemMsg, promptMsg, streamDeckButton, speechLanguage)
	aiBot.StreamDeckButtonWithHistory = streamDeckButtonWithHistory
	aiBot.StreamDeckButtonWithHistoryAndCopy = streamDeckButtonWithHistoryAndCopy
	aiBot.init()
	return aiBot
}
