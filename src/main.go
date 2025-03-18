package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/hvossi92/gollama/src/helpers"
	"github.com/hvossi92/gollama/src/services"
	"github.com/hvossi92/gollama/src/templates"
)

//go:embed static
var staticFs embed.FS

// Server struct to hold all services and templates
type Server struct {
	staticSubFS   fs.FS
	db            *services.DatabaseService
	ollamaService *services.OllamaService
}

// NewServer initializes and returns a new Server instance with all services set up.
func NewServer() (*Server, error) {
	// Create sub filesystem for static assets
	staticSubFS, err := fs.Sub(staticFs, "static")
	if err != nil {
		return nil, fmt.Errorf("failed to create sub filesystem: %w", err)
	}

	vectorDB, err := services.SetUpDatabaseService("database.db", false)
	if err != nil {
		return nil, fmt.Errorf("failed to set up VectorDB service: %w", err)
	}
	settings, err := vectorDB.GetSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}
	ollamaService := services.SetUpOllamaService(settings.URL, settings.LLM, settings.Embedding, staticFs)

	return &Server{
		staticSubFS:   staticSubFS,
		db:            vectorDB,
		ollamaService: ollamaService,
	}, nil
}

func main() {
	start := time.Now()
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
	defer server.db.Close()

	err = os.RemoveAll("./uploads")
	if err != nil {
		log.Fatalf("Failed to delete uploads directory: %v", err)
	}

	http.HandleFunc("/", server.fetchIndexPage)
	http.HandleFunc("/settings", server.fetchSettingsPage)
	http.HandleFunc("POST /settings", server.updateSettings)
	http.HandleFunc("POST /start", server.startAdventure)
	http.HandleFunc("POST /chat", server.fetchAiResponse)
	http.HandleFunc("GET /vector", server.GetVectors)
	http.HandleFunc("POST /vector", server.UploadVector)
	http.HandleFunc("PUT /settings", server.UpdateSettings)
	http.HandleFunc("GET /die", server.getDie)
	http.HandleFunc("/save-games", server.fetchSaveGamesPage)
	http.HandleFunc("GET /save-game/{id}", server.loadSaveGame)
	http.HandleFunc("POST /save-game", server.saveGame)
	http.HandleFunc("DELETE /save-game/{id}", server.deleteSaveGame)
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(server.staticSubFS))))

	fmt.Println("Server listening on port 2048")
	duration := time.Since(start)
	fmt.Printf("Server started in %s\n", duration)
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

	// settings, err := s.vectorDB.GetSettings()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// data := struct {
	// 	URL       string
	// 	LLM       string
	// 	Embedding string
	// }{
	// 	URL:       settings.URL,
	// 	LLM:       settings.LLM,
	// 	Embedding: settings.Embedding,
	// }

	err := templates.Index([]templ.Component{}).Render(r.Context(), w)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

func (s *Server) fetchSettingsPage(w http.ResponseWriter, r *http.Request) {
	templates.Settings().Render(r.Context(), w)
	if err := templates.Settings().Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) updateSettings(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Updating settings")
	settings := struct {
		URL       string
		LLM       string
		Embedding string
	}{
		URL:       r.FormValue("url"),
		LLM:       r.FormValue("llm"),
		Embedding: r.FormValue("embedding"),
	}
	fmt.Println(settings)
}

func (s *Server) startAdventure(w http.ResponseWriter, r *http.Request) {
	var err error
	fmt.Println("Starting adventure")
	// Setup random adventure using oracle
	myAdventure := services.ConsultAdventureOracle()
	fmt.Println(myAdventure)
	message := fmt.Sprintf("I am ready to start my adventure: \"{%s}\". What do I do?", myAdventure)
	aiResponse, err := s.ollamaService.AskLLM(message, s.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	firstMsgComponent := templates.Message(message, aiResponse)
	err = templates.ChatInterface([]templ.Component{firstMsgComponent}).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) fetchAiResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Println(services.ConsultD6Oracle())
	message := r.FormValue("message")

	var err error
	fmt.Println("Asking LLM")
	aiResponse, err := s.ollamaService.AskLLM(message, s.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	err = templates.Message(message, aiResponse).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) GetVectors(w http.ResponseWriter, r *http.Request) {
	text, err := s.db.ReadAllVectors()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.Message("Get all vectors", text).Render(r.Context(), w)
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

		err = s.db.InsertChunkAndEmbedding(chunk, embeddings) // Store chunk and embedding in DB
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	err := s.db.UpdateSettings(r.FormValue("url"), r.FormValue("llm"), r.FormValue("embedding"))
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

func (s *Server) getDie(w http.ResponseWriter, r *http.Request) {
	dieStr := r.URL.Query().Get("die")
	die, err := strconv.Atoi(dieStr)
	if err != nil {
		http.Error(w, "Invalid die value", http.StatusBadRequest)
		return
	}
	fmt.Print(rand.IntN(die))
}

func (s *Server) fetchSaveGamesPage(w http.ResponseWriter, r *http.Request) {
	saves, err := s.db.GetSaveGames()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.SaveGames(saves).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) loadSaveGame(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	fmt.Println(id)

	savedMessages, err := s.db.LoadSaveGame(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(savedMessages)

	s.ollamaService.LoadMessages(savedMessages)

	// firstMsgComponent := templates.Message(message, aiResponse)
	msgs := s.ollamaService.GetMessages()
	messages := []templ.Component{}
	for _, msg := range msgs {
		messages = append(messages, templates.Message(msg.Content, msg.Content))
	}
	fmt.Println(messages)
	err = templates.Index(messages).Render(r.Context(), w)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

func (s *Server) saveGame(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Convert messages to JSON string using the getter method
	content := s.ollamaService.GetMessages()

	err := s.db.SaveGame(name, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated table
	saves, err := s.db.GetSaveGames()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.SaveGames(saves).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) deleteSaveGame(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/save-game/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = s.db.DeleteSaveGame(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated table
	saves, err := s.db.GetSaveGames()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.SaveGames(saves).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
