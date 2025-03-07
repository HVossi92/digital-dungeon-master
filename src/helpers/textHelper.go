package helpers

import (
	"fmt"
	"strings"
)

func StripLAiResponse(text string) string {
	if containsTags(text, "think") {
		return stripTags(text, "think")
	}
	return text
}

func stripTags(input string, tag string) string {
	startTag := fmt.Sprintf("<%s>", tag)
	endTag := fmt.Sprintf("</%s>", tag)

	startIndex := strings.Index(input, startTag)
	endIndex := strings.LastIndex(input, endTag)

	if startIndex != -1 && endIndex != -1 && startIndex < endIndex {
		// <think>...</think> found and in correct order
		contentStart := startIndex + len(startTag)
		return input[:startIndex] + input[contentStart:endIndex] + input[endIndex+len(endTag):]
	}

	// <think>...</think> not found or not in correct format, return original string
	return input
}

func containsTags(input string, tag string) bool {
	startTag := fmt.Sprintf("<%s>", tag)
	endTag := fmt.Sprintf("</%s>", tag)

	startIndex := strings.Index(input, startTag)
	endIndex := strings.LastIndex(input, endTag)

	return startIndex != -1 && endIndex != -1 && startIndex < endIndex
}
