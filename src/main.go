package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hvossi92/gollama/src/helpers"
	"github.com/hvossi92/gollama/src/services"
)

//go:embed templates
var templatesFS embed.FS

//go:embed static
var staticFS embed.FS

// Server struct to hold all services and templates
type Server struct {
	templates     *template.Template
	staticSubFS   fs.FS
	vectorDB      *services.DatabaseService
	ollamaService *services.OllamaService
}

// NewServer initializes and returns a new Server instance with all services set up.
func NewServer() (*Server, error) {
	// Parse templates
	templates, err := template.ParseFS(templatesFS,
		"templates/*.html",
		"templates/**/*.html",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	// Create sub filesystem for static assets
	staticSubFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		return nil, fmt.Errorf("failed to create sub filesystem: %w", err)
	}

	vectorDB, err := services.SetUDatabaseService("database.db", false)
	if err != nil {
		return nil, fmt.Errorf("failed to set up VectorDB service: %w", err)
	}
	settings, err := vectorDB.GetSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}
	ollamaService := services.SetUpOllamaService(settings.URL, settings.LLM, settings.Embedding)

	return &Server{
		templates:     templates,
		staticSubFS:   staticSubFS,
		vectorDB:      vectorDB,
		ollamaService: ollamaService,
	}, nil
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
	defer server.vectorDB.Close() // Important to close VectorDB service when done

	err = os.RemoveAll("./uploads")
	if err != nil {
		log.Fatalf("Failed to delete uploads directory: %v", err)
	}

	http.HandleFunc("/", server.fetchIndexPage)
	http.HandleFunc("POST /chat", server.fetchAiResponse)
	http.HandleFunc("GET /vector", server.GetVectors)
	http.HandleFunc("POST /vector", server.UploadVector)
	http.HandleFunc("PUT /settings", server.UpdateSettings)
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(server.staticSubFS))))

	fmt.Println("Server listening on port 2048")
	err = http.ListenAndServe(":2048", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func (s *Server) fetchIndexPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	settings, err := s.vectorDB.GetSettings()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		URL       string
		LLM       string
		Embedding string
	}{
		URL:       settings.URL,
		LLM:       settings.LLM,
		Embedding: settings.Embedding,
	}

	err = s.templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

func (s *Server) fetchAiResponse(w http.ResponseWriter, r *http.Request) {
	message := r.FormValue("message")
	doUseRag := r.URL.Query().Get("use-rag") == "true"

	var err error
	fmt.Println("Asking LLM")
	aiResponse, err := s.ollamaService.AskLLM(message, doUseRag, s.vectorDB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	fmt.Printf("AI Response: %s", aiResponse)

	data := struct {
		UserMessage string
		AIResponse  string
	}{
		UserMessage: message,
		AIResponse:  aiResponse,
	}
	err = s.templates.ExecuteTemplate(w, "message.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) GetVectors(w http.ResponseWriter, r *http.Request) {
	text, err := s.vectorDB.ReadAllVectors()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := struct {
		UserMessage string
		AIResponse  string
	}{
		UserMessage: "Get all vectors",
		AIResponse:  text,
	}
	err = s.templates.ExecuteTemplate(w, "message.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) UploadVector(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("vectors")
	if text == "" {
		http.Error(w, "No data was provided", http.StatusBadRequest)
		return
	}
	chunkedText, err := helpers.ChunkText(strings.TrimSpace(text), 16, 4) // Chunk text first
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, chunk := range chunkedText { // Iterate through each text chunk

		embeddings, err := s.ollamaService.GetVectorEmbedding(chunk) // Get embedding for each chunk

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = s.vectorDB.InsertChunkAndEmbedding(chunk, embeddings) // Store chunk and embedding in DB
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	err := s.vectorDB.UpdateSettings(r.FormValue("url"), r.FormValue("llm"), r.FormValue("embedding"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte("Settings updated"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
