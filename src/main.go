package main

import (
	"context"
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
	"github.com/hvossi92/gollama/src/db"
	"github.com/hvossi92/gollama/src/services"
	"github.com/hvossi92/gollama/src/templates"
)

//go:embed static
var staticFs embed.FS

// Server struct to hold all services and templates
type Server struct {
	staticSubFS   fs.FS
	ollamaService *services.OllamaService
	queries       *db.Queries
	ctx           context.Context
}

// NewServer initializes and returns a new Server instance with all services set up.
func NewServer() (*Server, error) {
	// Create sub filesystem for static assets
	staticSubFS, err := fs.Sub(staticFs, "static")
	if err != nil {
		return nil, fmt.Errorf("failed to create sub filesystem: %w", err)
	}

	database := services.CreateDb("./db/database.db", false)
	// defer services.CloseDb(database)

	ctx := context.Background()
	// create tables
	if _, err := database.ExecContext(ctx, ddl); err != nil {
		log.Fatal(err)
	}
	queries := db.New(database)
	settings, err := queries.GetSettings(ctx)
	if err != nil {
		err = queries.InsertSettings(ctx, db.InsertSettingsParams{Url: "http://localhost:11434", Llm: "gemma3:4b"})
		if err != nil {
			panic(err)
		}
	}
	// settings, err := queries.GetSettings(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get settings: %w", err)
	// }
	ollamaService := services.SetUpOllamaService(settings.Url, settings.Llm, staticFs)

	return &Server{
		staticSubFS:   staticSubFS,
		ollamaService: ollamaService,
		queries:       queries,
		ctx:           ctx,
	}, nil
}

//go:embed db/schemas/schema.sql
var ddl string

func main() {
	start := time.Now()
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	err = os.RemoveAll("./uploads")
	if err != nil {
		log.Fatalf("Failed to delete uploads directory: %v", err)
	}

	// Register all routes
	server.RegisterRoutes()

	fmt.Println("Server listening on port 2048")
	duration := time.Since(start)
	fmt.Printf("Server started in %s\n", duration)
	err = http.ListenAndServe(":2048", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// RegisterRoutes sets up all the HTTP endpoints for the server
func (s *Server) RegisterRoutes() {
	http.HandleFunc("/", s.fetchIndexPage)
	http.HandleFunc("GET /settings", s.fetchSettingsPage)
	http.HandleFunc("PUT /settings", s.updateSettings)
	http.HandleFunc("POST /start", s.startAdventure)
	http.HandleFunc("POST /chat", s.fetchAiResponse)
	http.HandleFunc("GET /die", s.getDie)
	http.HandleFunc("/save-games", s.fetchSaveGamesPage)
	http.HandleFunc("GET /save-game/{id}", s.loadSaveGame)
	http.HandleFunc("POST /save-game", s.saveGame)
	http.HandleFunc("DELETE /save-game/{id}", s.deleteSaveGame)
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(s.staticSubFS))))
}

// API handler methods
func (s *Server) fetchIndexPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	err := templates.Index([]templ.Component{}).Render(r.Context(), w)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

func (s *Server) fetchSettingsPage(w http.ResponseWriter, r *http.Request) {
	settings, err := s.queries.GetSettings(s.ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := templates.Settings(settings).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) startAdventure(w http.ResponseWriter, r *http.Request) {
	var err error
	fmt.Println("Starting adventure")
	// Setup random adventure using oracle
	myAdventure := services.ConsultAdventureOracle()
	fmt.Println(myAdventure)
	message := fmt.Sprintf("I am ready to start my adventure: \"{%s}\". What do I do?", myAdventure)
	aiResponse, err := s.ollamaService.AskLLM(message)
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
	aiResponse, err := s.ollamaService.AskLLM(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	err = templates.Message(message, aiResponse).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) updateSettings(w http.ResponseWriter, r *http.Request) {
	settings := db.UpdateSettingsParams{
		Url: r.FormValue("url"),
		Llm: r.FormValue("llm"),
	}
	err := s.queries.UpdateSettings(s.ctx, settings)
	if err != nil {
		http.Error(w, "Invalid die value", http.StatusBadRequest)
	}
	s.fetchSettingsPage(w, r)
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
	saves, err := s.queries.GetSaveGames(s.ctx)
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

	savedMessages, err := s.queries.GetChatMessages(s.ctx, int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(savedMessages)

	chatMessages := make([]services.ChatMessage, len(savedMessages))
	for i, msg := range savedMessages {
		chatMessages[i] = services.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	s.ollamaService.LoadMessages(chatMessages)

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

	err := s.queries.InsertSaveGame(s.ctx, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert messages to JSON string using the getter method
	chatMessages := s.ollamaService.GetMessages()
	saveGame, err := s.queries.GetSaveGames(s.ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(saveGame) == 0 {
		http.Error(w, "No save games found", http.StatusInternalServerError)
		return
	}

	lastSaveGame := saveGame[len(saveGame)-1]

	for _, chatMessage := range chatMessages {
		params := db.InsertChatMessageParams{
			Role:       chatMessage.Role,
			Content:    chatMessage.Content,
			SaveGameID: lastSaveGame.ID,
		}
		err = s.queries.InsertChatMessage(s.ctx, params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Return updated table
	saves, err := s.queries.GetSaveGames(s.ctx)
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

	err = s.queries.DeleteSaveGame(s.ctx, int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated table
	saves, err := s.queries.GetSaveGames(s.ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.SaveGames(saves).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
