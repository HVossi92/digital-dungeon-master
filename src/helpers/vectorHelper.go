package helpers

import (
	"fmt"
	"regexp"
	"strings"
)

// chunkText chunks a string of text into smaller overlapping text chunks based on sentences.
//
// Parameters:
//
//	text:       The input string of text to chunk.
//	chunkSize:  The desired number of sentences per chunk.
//	chunkOverlap: The number of sentences to overlap between consecutive chunks.
//
// Returns:
//
//	[]string:  A slice of strings, where each string is a chunk of text (composed of sentences).
//	error:     An error if the input parameters are invalid.
func ChunkText(text string, chunkSize int, chunkOverlap int) ([]string, error) {
	if chunkSize <= 0 {
		return nil, fmt.Errorf("chunkSize must be greater than 0")
	}
	if chunkOverlap < 0 {
		return nil, fmt.Errorf("chunkOverlap must be non-negative")
	}
	if chunkOverlap >= chunkSize {
		return nil, fmt.Errorf("chunkOverlap must be less than chunkSize")
	}

	sentences := splitTextIntoSentences(text) // Split text into sentences
	if len(sentences) == 0 {
		return []string{}, nil // Return empty slice for empty input text
	}

	var chunks []string
	step := chunkSize - chunkOverlap

	for i := 0; i < len(sentences); i += step {
		end := i + chunkSize
		if end > len(sentences) {
			end = len(sentences) // Adjust end for the last chunk
		}
		chunkSentences := sentences[i:end]
		chunk := strings.Join(chunkSentences, " ") // Join sentences in the chunk back into a string
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

// SplitTextIntoSentences splits a text into sentences using a simple regex.
// Note: This is a basic sentence splitter and might not be perfect for all cases.
// For more robust sentence splitting, consider using NLP libraries.
func splitTextIntoSentences(text string) []string {
	// Use a regex to split by sentence-ending punctuation (. ! ?) followed by whitespace
	sentenceRegex := regexp.MustCompile(`(?P<sentence>[^.!?]+[.!?])\s+`) // Improved regex

	var sentences []string
	matches := sentenceRegex.FindAllStringSubmatch(text, -1)
	if matches == nil {
		// Handle case where no sentences are found using the regex (e.g., very short text)
		sentences = strings.Split(text, ". ") // Fallback to simple split if regex fails
		if len(sentences) <= 1 {              // If still only one or zero sentences after simple split
			sentences = strings.Split(text, "\n") // Try splitting by newline as a last resort
			if len(sentences) <= 1 {
				sentences = []string{text} // If all else fails, treat the whole text as one sentence
			}
		}
		return sentences

	}
	for _, match := range matches {
		sentences = append(sentences, strings.TrimSpace(match[1])) // match[1] is the captured sentence group
	}

	// Handle the last part of the text that might not end with sentence-ending punctuation
	lastSentence := sentenceRegex.ReplaceAllString(text, "")
	lastSentence = strings.TrimSpace(lastSentence)
	if lastSentence != "" {
		sentences = append(sentences, lastSentence)
	}

	return sentences
}
