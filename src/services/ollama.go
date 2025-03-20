package services

import (
	"embed"
	"fmt"
	"io/fs"
	"time"

	"github.com/hvossi92/gollama/src/helpers"
	"github.com/hvossi92/gollama/src/utils"
)

type OllamaService struct {
	chatEndpoint      string
	embeddingEndpoint string
	llm               string
	systemPrompt      string
	messages          []ChatMessage
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
	Model              string      `json:"model"`
	CreatedAt          string      `json:"created_at"`
	Message            ChatMessage `json:"message"` // <--- Changed to ChatMessageResponse struct
	DoneReason         string      `json:"done_reason"`
	Done               bool        `json:"done"`
	TotalDuration      int64       `json:"total_duration"`
	LoadDuration       int64       `json:"load_duration"`
	PromptEvalCount    int         `json:"prompt_eval_count"`
	PromptEvalDuration int64       `json:"prompt_eval_duration"`
	EvalCount          int         `json:"eval_count"`
	EvalDuration       int64       `json:"eval_duration"`
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

func SetUpOllamaService(url string, llm string, staticFs embed.FS) *OllamaService {
	filename := "static/systemPrompt.txt"
	content, err := fs.ReadFile(staticFs, filename)
	if err != nil {
		panic(err)
	}

	service := &OllamaService{
		chatEndpoint:      url + "/api/chat",
		embeddingEndpoint: url + "/api/embed",
		llm:               llm,
		systemPrompt:      string(content),
		messages:          make([]ChatMessage, 0),
	}

	systemPromptMsg := ChatMessage{
		Role:    "system",
		Content: service.systemPrompt,
	}
	service.messages = append(service.messages, systemPromptMsg)

	return service
}

func (s *OllamaService) AskLLM(question string) (string, error) {
	userMsg := ChatMessage{
		Role:    "user",
		Content: "{Remember asking for skill check, but only for the player and only when necessary} " + question,
	}
	s.messages = append(s.messages, userMsg)

	// 6. Make the Chat Request to Ollama
	request := ChatRequest{ // Use ChatRequest struct
		Model:    s.llm,
		Messages: s.messages,
		Stream:   false,
	}
	fmt.Println("Sending request to ollama")
	start := time.Now()
	chatResponse, err := utils.SendPostRequest[ChatRequest, ChatResponse](s.chatEndpoint, request) // Use ChatRequest and ChatResponse
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	fmt.Printf("Receiving response from '%s' took %s\n", s.llm, time.Since(start))

	responseTxt := helpers.StripLAiResponse(chatResponse.Message.Content)
	fmt.Println("Received response from ollama")

	aiResponseMsg := ChatMessage{
		Role:    "assistant",
		Content: responseTxt,
	}
	s.messages = append(s.messages, aiResponseMsg)

	return responseTxt, nil // Return response from LLM
}

func (s *OllamaService) GetMessages() []ChatMessage {
	return s.messages
}

func (s *OllamaService) LoadMessages(messages []ChatMessage) {
	s.messages = messages
}
