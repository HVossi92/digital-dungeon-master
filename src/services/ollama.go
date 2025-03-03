package services

import (
	"fmt"
	"strings"

	"github.com/hvossi92/gollama/src/utils"
)

type OllamaService struct {
	chatEndpoint      string
	embeddingEndpoint string
	llm               string
	embeddingModel    string
}

// ChatRequest struct to structure the request body
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct { // New struct for /api/chat response
	Model              string              `json:"model"`
	CreatedAt          string              `json:"created_at"`
	Message            ChatMessageResponse `json:"message"` // <--- Changed to ChatMessageResponse struct
	DoneReason         string              `json:"done_reason"`
	Done               bool                `json:"done"`
	TotalDuration      int64               `json:"total_duration"`
	LoadDuration       int64               `json:"load_duration"`
	PromptEvalCount    int                 `json:"prompt_eval_count"`
	PromptEvalDuration int64               `json:"prompt_eval_duration"`
	EvalCount          int                 `json:"eval_count"`
	EvalDuration       int64               `json:"eval_duration"`
}

type ChatMessageResponse struct { // Struct for the nested "message" object
	Role    string `json:"role"`
	Content string `json:"content"` // <--- The text response is here
}

type EmbeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbeddingResponse struct {
	Model             string      `json:"model"`
	Embeddings        [][]float32 `json:"embeddings"`
	Total_duration    int
	Load_duration     int
	Prompt_eval_count int
}

// Point struct for x, y coordinates
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// SetUpVectorDBService creates and initializes a new VectorDBService.
func SetUpOllamaService(url string, llm string, embedding string) *OllamaService {
	return &OllamaService{chatEndpoint: url + "/api/chat", embeddingEndpoint: url + "/api/embed", llm: llm, embeddingModel: embedding}
}

var questionSystemPrompt = `
You are a helpful assistant with access to a knowlege base, tasked with answering questions about general knowledge, but also specific to the provided knowledge base.

Answer the question in a very concise manner. Use an unbiased and journalistic tone. Do not repeat text. Don't make anything up. If you are not sure about something, just say that you don't know.
{{- /* Stop here if no context is provided. The rest below is for handling contexts. */ -}}
{{- if . -}}
If possible, answer the question solely based on the provided search results from the knowledge base. If the search results from the knowledge base are not relevant to the question at hand, try to answer the question based on general knowledge. But do not make anything up.

Anything between the following 'context' XML blocks is retrieved from the knowledge base, not part of the conversation with the user. The bullet points are ordered by relevance, so the first one is the most relevant.

<context>
    {{- if . -}}
    {{- range $context := .}}
    - {{.}}{{end}}
    {{- end}}
</context>
{{- end -}}

Don't mention the knowledge base, context or search results in your answer.
`

func (s *OllamaService) AskLLM(question string, useVectorDb bool, vectorService *DatabaseService) (string, error) {
	var messages []ChatMessage

	if useVectorDb {
		// 1. Embed the question to find relevant chunks
		questionEmbedding, err := s.GetVectorEmbedding(question)
		if err != nil {
			return "", fmt.Errorf("failed to embed question: %w", err)
		}

		// 2. Query vector DB to find similar chunks
		similarItems, err := vectorService.FindSimilarVectors(questionEmbedding)
		if err != nil {
			return "", fmt.Errorf("failed to find similar vectors: %w", err)
		}

		// 3. Construct context from retrieved chunks
		context := ""
		if len(similarItems) > 0 {
			contextBuilder := strings.Builder{}
			contextBuilder.WriteString("Context:\n")
			for _, item := range similarItems {
				contextBuilder.WriteString(item.Text)
				contextBuilder.WriteString("\n---\n") // Separator between chunks
			}
			context = contextBuilder.String()
		} else {
			context = "No relevant context found in the database.\n"
		}

		// 4. Create prompt with context and question
		messages = []ChatMessage{
			{
				Role:    "system",
				Content: questionSystemPrompt,
			}, {
				Role:    "user",
				Content: "<context>" + context + "</context>" + "\nQuestion: " + question, // Combine context and question
			},
		}
	} else {
		// 5. If not using vector DB, use a simple prompt with just the question
		messages = []ChatMessage{
			{
				Role:    "system",
				Content: questionSystemPrompt,
			}, {
				Role:    "user",
				Content: "Question: " + question,
			},
		}
	}

	// 6. Make the Chat Request to Ollama
	request := ChatRequest{ // Use ChatRequest struct
		Model:    s.llm,
		Messages: messages,
		Stream:   false,
	}
	chatResponse, err := utils.SendPostRequest[ChatRequest, ChatResponse](s.chatEndpoint, request) // Use ChatRequest and ChatResponse
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return chatResponse.Message.Content, nil // Return response from LLM
}

func (s *OllamaService) GetVectorEmbedding(text string) ([]float32, error) {
	request := EmbeddingRequest{
		Model: s.embeddingModel,
		Input: text,
	}

	fmt.Println("Generating vector embeddings", s.embeddingEndpoint)
	ollamaResponse, err := utils.SendPostRequest[EmbeddingRequest, EmbeddingResponse](s.embeddingEndpoint, request)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return ollamaResponse.Embeddings[0], nil // Return ollamaResponse.Message.Content
}
