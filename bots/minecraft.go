package bots

import (
	"OpenAITest/model"
	sd "OpenAITest/streamdeck"
	"OpenAITest/utils"
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/micmonay/keybd_event"
	openai "github.com/sashabaranov/go-openai"
	"golang.design/x/clipboard"
	"log"
	"strings"
	"time"
)

func TypeMultipleMinecraftCommands(commands string, kb *keybd_event.KeyBonding) error {
	for _, command := range strings.Split(commands, "\n") {
		if command == "" {
			continue
		}
		err := TypeMinecraftCommand(command, kb)
		if err != nil {
			return err
		}
		time.Sleep(20 * time.Millisecond)
	}
	return nil
}

func TypeMinecraftCommand(command string, kb *keybd_event.KeyBonding) error {
	time.Sleep(50 * time.Millisecond)
	kb.SetKeys(keybd_event.VK_T)
	err := kb.Launching()
	if err != nil {
		fmt.Printf("Error launching keyboard: %v\n", err)
		return err
	}
	time.Sleep(50 * time.Millisecond)
	kb.SetKeys(keybd_event.VK_BACKSPACE)
	err = kb.Launching()
	if err != nil {
		fmt.Printf("Error launching keyboard: %v\n", err)
		return err
	}

	if !strings.HasPrefix(command, "/") {
		command = "/" + command
	}
	clipboard.Write(clipboard.FmtText, []byte(command))
	time.Sleep(50 * time.Millisecond)
	kb.HasCTRL(true)
	kb.SetKeys(keybd_event.VK_V)
	err = kb.Launching()
	if err != nil {
		fmt.Printf("Error launching keyboard: %v\n", err)
		return err
	}
	time.Sleep(50 * time.Millisecond)
	kb.HasCTRL(false)
	kb.SetKeys(keybd_event.VK_ENTER)
	err = kb.Launching()
	if err != nil {
		fmt.Printf("Error launching keyboard: %v\n", err)
		return err
	}
	return nil
}

var (
	quitChannel chan bool
	finished    chan bool
	isRecording bool
)

func init() {
	quitChannel = make(chan bool)
	finished = make(chan bool)
	isRecording = false
}

func EvaluateMinecraftGptResponseStrings(input []string, chatContent model.ChatContent, client *openai.Client, speech *htgotts.Speech, kb *keybd_event.KeyBonding) {
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
	fmt.Printf("Code blocks: %v\n", codeBlocks)
	if len(codeBlocks) == 0 {
		log.Println("No code block found")
	}
	fmt.Println("Code blocks found: ", len(codeBlocks))

	commandsToExecute := strings.Join(codeBlocks, "\n")
	err = TypeMultipleMinecraftCommands(commandsToExecute, kb)
	if err != nil {
		fmt.Printf("Error typing command: %v\n", err)
	}
}

func InitMinecraftGPTBot(client *openai.Client, device sd.DeviceWrapper, properties map[string]string,
	streamdeckHandler sd.IStreamdeckHandler, speech *htgotts.Speech, kb *keybd_event.KeyBonding, button uint8) {
	minecraftSystemMsg := properties["minecraftSystemMsg"]
	minecraftPromptMsg := properties["minecraftPromptMsg"]

	minecraftChatContent := model.ChatContent{
		SystemMsg: minecraftSystemMsg,
		PromptMsg: minecraftPromptMsg,
	}

	err := sd.SetStreamdeckButtonText(device, button, "Minecraft")
	if err != nil {
		log.Fatal(err)
	}

	streamdeckHandler.AddOnPressHandler(int(button), func() error {
		go func() {
			isRecording = true
			utils.RecordAndSaveAudioAsMp3("minecraft.wav", quitChannel, finished)
		}()
		return nil
	})
	streamdeckHandler.AddOnReleaseHandler(int(button), func() error {
		if isRecording {
			quitChannel <- true
			<-finished
			fmt.Println("OnRealease: Stopping recording...")
			isRecording = false
			transcription, err := utils.ParseMp3ToText("minecraft.wav", client)
			if err != nil {
				fmt.Printf("Error parsing mp3 to text: %s\n", err)
				return nil
			}
			EvaluateMinecraftGptResponseStrings([]string{transcription}, minecraftChatContent, client, speech, kb)
		}
		return nil
	})
}
