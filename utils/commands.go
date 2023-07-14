package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/micmonay/keybd_event"
	"golang.design/x/clipboard"
	"os"
	"os/exec"
	"strconv"
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

func SplitStringByDot(input string) []string {
	return strings.Split(input, ".")
}

func SplitStringByComma(input string) []string {
	return strings.Split(input, ",")
}

func SplitStringIntoLogicalChunks(input string, chunkSize int) []string {
	splitStringByDot := SplitStringByDot(input)
	var result []string
	//currentChunk := ""

	for _, sentence := range splitStringByDot {
		if len(sentence) < chunkSize {
			result = append(result, sentence)
			continue
		}
		commaSplit := SplitStringByComma(sentence)
		for _, commaSentence := range commaSplit {
			if len(commaSentence) < chunkSize {
				result = append(result, commaSentence)
				continue
			}
			chunks := SplitStringIntoWordChunks(commaSentence, chunkSize)
			result = append(result, chunks...)
		}
	}
	return result
}

func SplitStringIntoWordChunks(input string, chunkSize int) []string {
	var result []string
	words := strings.Fields(input)
	currentChunk := ""

	for _, word := range words {
		if len(currentChunk)+len(word) > chunkSize {
			result = append(result, currentChunk)
			currentChunk = word
		} else {
			if currentChunk == "" {
				currentChunk = word
			} else {
				currentChunk += " " + word
			}
		}
	}

	if currentChunk != "" {
		result = append(result, currentChunk)
	}

	return result
}

func ReadNumber(commands []string) int {
	fmt.Println("Available Command: ")
	fmt.Printf("0. Cancel\n")
	for i, command := range commands {
		fmt.Printf("%d. %s\n", i+1, command)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter a number: ")
	scanner.Scan()
	input := scanner.Text()
	num, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("%s is not a valid number.\n", input)
		return -1
	}
	return num
}

func RunAndWaitForCommandSelection(scanner *bufio.Scanner, commands []string) int {
	fmt.Println("Available Command: ")
	for i, command := range commands {
		fmt.Printf("0. Cancel\n")
		fmt.Printf("%d. %s\n", i, command)
	}

	for {
		fmt.Print("Enter Command Number: \n")
		scanner.Scan()
		text := scanner.Text()
		fmt.Println("Inputs is: ", text)

		trimmedText := strings.TrimSpace(text)
		if len(trimmedText) != 0 {
			atoi, err := strconv.Atoi(trimmedText)
			if err != nil {
				fmt.Println("Invalid Inputs")
				continue
			}
			return atoi
		}
	}
}

func printCommandOutput(command string) error {
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func RunTerminalCommand(command string) (string, error) {
	commandWithoutNewlines := strings.ReplaceAll(command, "\n", " ")
	cmd := exec.Command("/bin/sh", "-c", commandWithoutNewlines)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	if out.Len() == 0 && errOut.Len() != 0 {
		return "", fmt.Errorf(errOut.String())
	}
	return out.String(), nil
}
