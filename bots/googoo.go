package bots

import (
	"OpenAITest/model"
	"OpenAITest/utils"
	"bufio"
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/sashabaranov/go-openai"
	"log"
	"strings"
)

//1. use googler --noprompt to search for the question, download the first result and write the content
//2. Remove al useless content
//3. Use the content to answer the question

var googooCompletionHistory []openai.ChatCompletionMessage

func EvaluateGoogooGptResponseStrings(input []string, withHistory bool, scanner *bufio.Scanner, chatContent model.ChatContent, client *openai.Client) {
	joinedRequestMessage := strings.Join(input, "\n")
	if withHistory {
		chatContent.HistoryMessages = googooCompletionHistory
	} else {
		googooCompletionHistory = []openai.ChatCompletionMessage{}
	}
	answer, err := utils.SendChatRequest(joinedRequestMessage, chatContent, client)
	if err != nil {
		log.Fatal(err)
	}

	googooCompletionHistory = append(googooCompletionHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: joinedRequestMessage,
	})
	googooCompletionHistory = append(googooCompletionHistory, openai.ChatCompletionMessage{
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
			commandReturnValue, err = utils.RunTerminalCommand(codeBlocks[commandToExecute+1])
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

//func InitGoogooGPTBot(client *openai.Client, device *streamdeck.Device, properties map[string]string,
//	streamdeckHandler streamdeck2.IStreamdeckHandler, scanner *bufio.Scanner, buttonWithoutHistory int16, buttonWithHistory int16) *model.ChatContent {
//
//	googooSystemMsg := properties["googooSystemMsg"]
//	googooPromptMsg := properties["googooPromptMsg"]
//	googooChatContent := model.ChatContent{
//		SystemMsg:       googooSystemMsg,
//		PromptMsg:       googooPromptMsg,
//		HistoryMessages: []openai.ChatCompletionMessage{},
//	}
//
//	if buttonWithoutHistory >= 0 {
//		err := streamdeckHandler.AddButtonText(int(buttonWithoutHistory), "Googoo")
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		streamdeckHandler.AddOnPressHandler(int(buttonWithoutHistory), func() error {
//			go func() {
//				isRecording = true
//				utils.RecordAndSaveAudioAsMp3("Googoo.wav", quitChannel, finished)
//			}()
//			return nil
//		})
//
//		streamdeckHandler.AddOnReleaseHandler(int(buttonWithoutHistory), func() error {
//			if isRecording {
//				quitChannel <- true
//				<-finished
//				isRecording = false
//				transcription, err := utils.ParseMp3ToText("Googoo.wav", client)
//				if err != nil {
//					fmt.Printf("Error parsing mp3 to text: %s\n", err)
//					return nil
//				}
//				EvaluateGoogooGptResponseStrings([]string{transcription}, false, scanner, googooChatContent, client)
//			}
//			return nil
//		})
//	}
//
//	if buttonWithHistory >= 0 {
//		err := streamdeckHandler.AddButtonText(int(buttonWithHistory), "HGoogoo")
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		streamdeckHandler.AddOnPressHandler(int(buttonWithHistory), func() error {
//			go func() {
//				isRecording = true
//				utils.RecordAndSaveAudioAsMp3("Googoo.wav", quitChannel, finished)
//			}()
//			return nil
//		})
//
//		streamdeckHandler.AddOnReleaseHandler(int(buttonWithHistory), func() error {
//			if isRecording {
//				quitChannel <- true
//				<-finished
//				isRecording = false
//				transcription, err := utils.ParseMp3ToText("Googoo.wav", client)
//				if err != nil {
//					fmt.Printf("Error parsing mp3 to text: %s\n", err)
//					return nil
//				}
//				EvaluateGoogooGptResponseStrings([]string{transcription}, true, scanner, googooChatContent, client)
//			}
//			return nil
//		})
//	}
//	return &googooChatContent
//}
