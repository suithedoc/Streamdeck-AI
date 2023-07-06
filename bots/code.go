package bots

import (
	"OpenAITest/model"
	streamdeck2 "OpenAITest/streamdeck"
	"OpenAITest/utils"
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/micmonay/keybd_event"
	"github.com/sashabaranov/go-openai"
	"golang.design/x/clipboard"
	"log"
	"strings"
	"time"
)

func TypeCodeCommand(command string, kb *keybd_event.KeyBonding) error {
	clipboard.Write(clipboard.FmtText, []byte(command))
	time.Sleep(100 * time.Millisecond)
	kb.HasCTRL(true)
	kb.SetKeys(keybd_event.VK_V)
	err := kb.Launching()
	if err != nil {
		fmt.Printf("Error launching keyboard: %v\n", err)
		return err
	}
	kb.HasCTRL(false)
	time.Sleep(20 * time.Millisecond)
	return nil
}

func init() {
	quitChannel = make(chan bool)
	isRecording = false
}

func EvaluateCodeGptResponseStrings(input []string, chatContent model.ChatContent, client *openai.Client, speech *htgotts.Speech, kb *keybd_event.KeyBonding) {
	answer, err := utils.SendChatRequest(strings.Join(input, "\n"), chatContent, client)
	if err != nil {
		log.Fatal(err)
	}
	result := markdown.Render(answer, 80, 6)
	fmt.Println(string(result))
	codeBlocks := utils.ExtractCodeBlockFromMarkdown(answer)
	if len(codeBlocks) == 0 {
		codeBlocks = utils.ExtractCodeBlockFromMarkdownWithOneBacktick(answer)
	}
	if len(codeBlocks) == 0 {
		log.Println("No code block found")
		fmt.Println("No Code blocks found in this: ", answer)
	}
	fmt.Println("Code blocks found: ", len(codeBlocks))

	allCodeBlocks := strings.Join(codeBlocks, "\n")
	err = TypeCodeCommand(allCodeBlocks, kb)
	if err != nil {
		fmt.Printf("Error typing command: %v\n", err)
	}
}

func InitCodeGPTBot(client *openai.Client, device streamdeck2.DeviceWrapper, properties map[string]string,
	streamdeckHandler streamdeck2.IStreamdeckHandler, speech *htgotts.Speech, kb *keybd_event.KeyBonding, button uint8) {

	codeChatContent := model.ChatContent{
		SystemMsg: properties["codeSystemMsg"],
		PromptMsg: properties["codePromptMsg"],
	}

	err := streamdeckHandler.AddButtonText(int(button), "Code")
	if err != nil {
		log.Fatal(err)
	}

	streamdeckHandler.AddOnPressHandler(int(button), func() error {
		go func() {
			isRecording = true
			utils.RecordAndSaveAudioAsMp3("Code.wav", quitChannel, finished)
		}()
		return nil
	})
	streamdeckHandler.AddOnReleaseHandler(streamdeckHandler.ReverseTraverseButtonId(int(button)), func() error {
		if isRecording {
			quitChannel <- true
			<-finished
			isRecording = false
			transcription, err := utils.ParseMp3ToText("Code.wav", client)
			if err != nil {
				fmt.Printf("Error parsing mp3 to text: %s\n", err)
				return nil
			}
			EvaluateCodeGptResponseStrings([]string{transcription}, codeChatContent, client, speech, kb)
		}
		return nil
	})
}
