package services

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/hvossi92/gollama/src/helpers"
	"github.com/hvossi92/gollama/src/utils"
)

type OllamaService struct {
	chatEndpoint      string
	embeddingEndpoint string
	llm               string
	embeddingModel    string
	systemPrompt      string
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

// SetUpVectorDBService creates and initializes a new VectorDBService.
func SetUpOllamaService(url string, llm string, embedding string, staticFs embed.FS) *OllamaService {
	filename := "static/systemPrompt.txt" // Replace with your text file name
	content, err := fs.ReadFile(staticFs, filename)
	if err != nil {
		panic(err)
	}

	return &OllamaService{chatEndpoint: url + "/api/chat", embeddingEndpoint: url + "/api/embed", llm: llm, embeddingModel: embedding, systemPrompt: string(content)}
}

func (s *OllamaService) AskLLM(question string, vectorService *DatabaseService) (string, error) {
	var messages []ChatMessage

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
		contextBuilder.WriteString("D&D Rule References:\n")
		for _, item := range similarItems {
			contextBuilder.WriteString(fmt.Sprintf("- %s\n", item.Text))
		}
		context = contextBuilder.String()
	} else {
		context = "No relevant context found in the database.\n"
	}

	// fmt.Println(context)
	// 4. Create prompt with context and question
	messages = []ChatMessage{
		{
			Role:    "system",
			Content: s.systemPrompt,
		}, {
			Role: "user",
			Content: fmt.Sprintf(
				"Rules Reference:\n%s\nPlayer Action: %s\nDM Response:",
				context,
				question,
			),
		},
	}

	// 6. Make the Chat Request to Ollama
	request := ChatRequest{ // Use ChatRequest struct
		Model:    s.llm,
		Messages: messages,
		Stream:   false,
	}
	fmt.Println("Sending request to ollama")
	chatResponse, err := utils.SendPostRequest[ChatRequest, ChatResponse](s.chatEndpoint, request) // Use ChatRequest and ChatResponse
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	fmt.Println("Received response from ollama")

	responseTxt := helpers.StripLAiResponse(chatResponse.Message.Content)

	return responseTxt, nil // Return response from LLM
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
