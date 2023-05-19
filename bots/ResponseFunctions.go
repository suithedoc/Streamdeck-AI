package bots

import (
	"OpenAITest/utils"
	"fmt"
	"log"
)

func SpeakResultFunc(bot *AiBot, answer string) error {
	chunks := utils.SplitStringIntoLogicalChunks(answer, 100)
	for _, chunk := range chunks {
		err := bot.speech.Speak(chunk)
		if err != nil {
			return fmt.Errorf("speaking answer: %v", err)
		}
	}
	return nil
}

func ExecuteCommandResultFunc(bot *AiBot, answer string) error {
	codeBlocks := utils.ExtractCodeBlockFromMarkdown(answer)
	if len(codeBlocks) == 0 {
		codeBlocks = utils.ExtractCodeBlockFromMarkdownWithOneBacktick(answer)
	}
	if len(codeBlocks) == 0 {
		log.Println("No code block found")
		log.Printf("Answer: \n%v", answer)
	}

	if len(codeBlocks) > 1 {
		//commandToExecute := utils.RunAndWaitForCommandSelection(scanner, codeBlocks)
		commandToExecute := utils.ReadNumber(codeBlocks)
		if commandToExecute <= 0 {
			return fmt.Errorf("invalid command: %v", commandToExecute)
		}
		if commandToExecute > len(codeBlocks) {
			return fmt.Errorf("invalid command: %v", commandToExecute)
		}
		fmt.Printf("Executing command: %v\n", commandToExecute)
		go func() {
			commandReturnValue, err := utils.RunTerminalCommand(codeBlocks[commandToExecute-1])
			if err != nil {
				fmt.Printf("Error running command: %v\n", err)
				return
			}
			fmt.Printf("Command return value: %v\n", commandReturnValue)

		}()

	} else if len(codeBlocks) == 1 {
		go func() {
			commandReturnValue, err := utils.RunTerminalCommand(codeBlocks[0])
			if err != nil {
				fmt.Printf("Error running command: %v\n", err)
				return
			}
			fmt.Printf("Command return value: %v\n", commandReturnValue)
		}()
	}
	return nil
}
