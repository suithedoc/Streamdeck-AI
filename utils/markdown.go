package utils

import (
	"fmt"
	"regexp"
)

func ExtractCodeBlockFromMarkdownWithOneBacktick(markdownText string) []string {
	// Define the regular expression to match code blocks within backticks
	re := regexp.MustCompile("`([^`]+?)`")

	// Find all matches in the input string
	matches := re.FindAllStringSubmatch(markdownText, -1)
	// Create a slice to store the code blocks
	codeBlocks := make([]string, len(matches))

	// Extract the code blocks from the matches and store them in the slice
	for i, match := range matches {
		codeBlocks[i] = match[1]
	}
	fmt.Printf("ExtractCodeBlockFromMarkdownWithOneBacktick: Found %d code blocks\n", len(matches))
	return codeBlocks
}

func ExtractCodeBlockFromMarkdown(markdownText string) []string {
	// Regular expression to match code blocks
	re := regexp.MustCompile("```(.*)\n([\\s\\S]*?)\n```")

	// Find all matches
	matches := re.FindAllStringSubmatch(markdownText, -1)

	var codeBlocks []string
	fmt.Printf("ExtractCodeBlockFromMarkdown: Found %d code blocks\n", len(matches))
	// Loop through matches and print code blocks
	for _, match := range matches {
		language := match[1]
		codeBlock := match[2]
		codeBlocks = append(codeBlocks, codeBlock)
		fmt.Printf("Language: %s\nCode block:\n%s\n\n", language, codeBlock)
	}
	return codeBlocks
}
