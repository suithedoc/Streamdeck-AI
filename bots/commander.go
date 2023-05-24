package bots

import (
	"OpenAITest/model"
	sd "OpenAITest/streamdeck"
	"OpenAITest/utils"
	"bufio"
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/sashabaranov/go-openai"
	"log"
	"strings"
)

var commanderCompletionHistory []openai.ChatCompletionMessage

func EvaluateCommanderGptResponseStrings(input []string, withHistory bool, scanner *bufio.Scanner, chatContent model.ChatContent, client *openai.Client) {
	joinedRequestMessage := strings.Join(input, "\n")
	if withHistory {
		chatContent.HistoryMessages = commanderCompletionHistory
	} else {
		commanderCompletionHistory = []openai.ChatCompletionMessage{}
	}
	answer, err := utils.SendChatRequest(joinedRequestMessage, chatContent, client)
	if err != nil {
		log.Fatal(err)
	}

	commanderCompletionHistory = append(commanderCompletionHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: joinedRequestMessage,
	})
	commanderCompletionHistory = append(commanderCompletionHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: answer,
	})

	result := markdown.Render(answer, 80, 6)
	fmt.Println(string(result))
	codeBlocks := utils.ExtractCodeBlockFromMarkdown(answer)
	if len(codeBlocks) == 0 {
		codeBlocks = utils.ExtractCodeBlockFromMarkdownWithOneBacktick(answer)
	}
	if len(codeBlocks) == 0 {
		log.Println("No code block found")
		log.Printf("Answer: \n%v", answer)
	}

	var commandReturnValue string
	if len(codeBlocks) > 1 {
		//commandToExecute := utils.RunAndWaitForCommandSelection(scanner, codeBlocks)
		commandToExecute := utils.ReadNumber(codeBlocks)
		if commandToExecute <= 0 {
			fmt.Println("Invalid command: ", commandToExecute)
			return
		}
		if commandToExecute > len(codeBlocks) {
			fmt.Println("Invalid command: ", commandToExecute)
			return
		}
		fmt.Printf("Executing command: %v\n", commandToExecute)
		go func() {
			commandReturnValue, err = utils.RunTerminalCommand(codeBlocks[commandToExecute-1])
			if err != nil {
				fmt.Printf("Error running command: %v\n", err)
				return
			}
			fmt.Printf("Command return value: %v\n", commandReturnValue)

		}()

	} else if len(codeBlocks) == 1 {
		go func() {
			commandReturnValue, err = utils.RunTerminalCommand(codeBlocks[0])
			if err != nil {
				fmt.Printf("Error running command: %v\n", err)
				return
			}
			fmt.Printf("Command return value: %v\n", commandReturnValue)

		}()
	}
}

func InitCommanderGPTBot(client *openai.Client, device sd.DeviceWrapper, properties map[string]string,
	streamdeckHandler sd.IStreamdeckHandler, scanner *bufio.Scanner, buttonWithoutHistory int16, buttonWithHistory int16) *model.ChatContent {

	commanderSystemMsg := properties["commanderSystemMsg"]
	commanderPromptMsg := properties["commanderPromptMsg"]
	commanderChatContent := model.ChatContent{
		SystemMsg:       commanderSystemMsg,
		PromptMsg:       commanderPromptMsg,
		HistoryMessages: []openai.ChatCompletionMessage{},
	}

	if buttonWithoutHistory >= 0 {
		err := sd.SetStreamdeckButtonText(device, uint8(buttonWithoutHistory), "Commander")
		if err != nil {
			log.Fatal(err)
		}

		streamdeckHandler.AddOnPressHandler(int(buttonWithoutHistory), func() error {
			go func() {
				isRecording = true
				utils.RecordAndSaveAudioAsMp3("audio.wav", quitChannel, finished)
			}()
			return nil
		})

		streamdeckHandler.AddOnReleaseHandler(int(buttonWithoutHistory), func() error {
			if isRecording {
				quitChannel <- true
				<-finished
				isRecording = false
				transcription, err := utils.ParseMp3ToText("audio.wav", client)
				if err != nil {
					fmt.Printf("Error parsing mp3 to text: %s\n", err)
					return nil
				}
				EvaluateCommanderGptResponseStrings([]string{transcription}, false, scanner, commanderChatContent, client)
			}
			return nil
		})
	}

	if buttonWithHistory >= 0 {
		err := sd.SetStreamdeckButtonText(device, uint8(buttonWithHistory), "HCommander")
		if err != nil {
			log.Fatal(err)
		}

		streamdeckHandler.AddOnPressHandler(int(buttonWithHistory), func() error {
			go func() {
				isRecording = true
				utils.RecordAndSaveAudioAsMp3("audio.wav", quitChannel, finished)
			}()
			return nil
		})

		streamdeckHandler.AddOnReleaseHandler(int(buttonWithHistory), func() error {
			if isRecording {
				quitChannel <- true
				<-finished
				isRecording = false
				transcription, err := utils.ParseMp3ToText("audio.wav", client)
				if err != nil {
					fmt.Printf("Error parsing mp3 to text: %s\n", err)
					return nil
				}
				EvaluateCommanderGptResponseStrings([]string{transcription}, true, scanner, commanderChatContent, client)
			}
			return nil
		})
	}
	return &commanderChatContent
}
