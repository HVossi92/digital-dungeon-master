package services

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/tursodatabase/go-libsql"
)

// DatabaseService represents the service responsible for the vector database.
type DatabaseService struct {
	db *sqlx.DB
}

type VectorItem struct {
	Id        int
	Text      string
	Embedding []byte
}

type Settings struct {
	URL       string
	LLM       string
	Embedding string
}

// SetUpDatabaseService creates and initializes a new VectorDBService.
func SetUpDatabaseService(dbPath string, overwrite bool) (*DatabaseService, error) {
	if overwrite {
		log.Println("Overwriting existing database (if it exists)")
		if err := os.Remove(dbPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("failed to remove existing database: %w", err)
		}
	}

	// Create VectorService instance first
	databaseService := &DatabaseService{db: nil}
	db, err := databaseService.createDb(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	databaseService.db = db

	// Then call ensureVectorTableExists on the instance
	if err := databaseService.createRulesVectorTable(); err != nil {
		db.Close() // Close the connection if table creation fails
		return nil, fmt.Errorf("failed to ensure vector table exists: %w", err)
	}

	if err := databaseService.createSettingsTable(); err != nil {
		db.Close() // Close the connection if table creation fails
		return nil, fmt.Errorf("failed to ensure settings table exists: %w", err)
	}

	return databaseService, nil
}

// Close closes the database connection.  Good practice to add a Close method.
func (s *DatabaseService) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *DatabaseService) createDb(dbPath string) (*sqlx.DB, error) {
	// Connect to embedded libSQL
	db, err := sqlx.Open("libsql", "file:"+dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Connection failed:", err)
	}

	log.Println("Connected to local libSQL database!")
	return db, nil
}

// EnsureVectorTableExists checks if the vector table exists and creates it if not.
func (s *DatabaseService) createRulesVectorTable() error {
	_, err := s.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS rules (id INTEGER PRIMARY KEY, text TEXT, embedding F32_BLOB(%d))", 768))
	if err != nil {
		return err
	}

	return nil
}

func (s *DatabaseService) createSettingsTable() error {
	// Create table if not exists
	s.db.MustExec(`CREATE TABLE IF NOT EXISTS settings (
		id INTEGER PRIMARY KEY, 
		url TEXT, 
		llm TEXT, 
		embedding_model TEXT
	) STRICT`)

	// Check if table is empty
	var count int
	err := s.db.Get(&count, "SELECT COUNT(*) FROM settings")
	if err != nil {
		return fmt.Errorf("failed to check settings table: %w", err)
	}

	// Insert default values only if table is empty
	if count == 0 {
		_, err = s.db.Exec(`INSERT INTO settings (url, llm, embedding_model) 
			VALUES (?, ?, ?)`,
			"http://192.168.178.105:11434",
			"llama3.1:8b-instruct-q8_0",
			"nomic-embed-text:latest")
		if err != nil {
			return fmt.Errorf("failed to insert default settings: %w", err)
		}
	}

	return nil
}

// InsertChunkAndEmbedding saves a text chunk and its embedding to the SQLite vector database.
func (s *DatabaseService) InsertChunkAndEmbedding(chunk string, embedding []float32) error {
	var sb strings.Builder
	sb.WriteByte('[')
	for i, v := range embedding {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(strconv.FormatFloat(float64(v), 'f', 6, 32))
	}
	sb.WriteByte(']')
	vectorStr := sb.String()

	_, err := s.db.Exec(
		`INSERT INTO rules (text, embedding) 
         VALUES (?, vector32(?))`,
		chunk,
		vectorStr,
	)
	if err != nil {
		return err
	}
	return nil
}

// ReadAllVectors reads all data from the vector DB table.
func (s *DatabaseService) ReadAllVectors() (string, error) {
	var vectorItems []VectorItem

	err := s.db.Select(&vectorItems, `
		SELECT id, text, vector_extract(embedding) 
		FROM rules
		ORDER BY id ASC`)
	if err != nil {
		return "", fmt.Errorf("query failed: %w", err)
	}

	var builder strings.Builder
	for _, vec := range vectorItems {
		display := vec.Text
		if len(vec.Text) > 40 {
			display = vec.Text[:40-3] + "..."
		}
		builder.WriteString(fmt.Sprintf("%-4d - %-40s | ", vec.Id, display))
	}
	return builder.String(), nil
}

// FindSimilarVectors queries the vector DB for vectors similar to the given embedding.
func (s *DatabaseService) FindSimilarVectors(queryEmbedding []float32) ([]VectorItem, error) {
	var sb strings.Builder
	sb.WriteByte('[')
	for i, v := range queryEmbedding {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(strconv.FormatFloat(float64(v), 'f', 6, 32))
	}
	sb.WriteByte(']')
	vectorStr := sb.String()

	var similarItems []VectorItem
	// TODO: Might need a distance field in the struct, or alternatively do not return the distance and only use it for ordering (since atm I only use it inside the db)
	err := s.db.Select(&similarItems,
		`SELECT text, vector_extract(embedding), vector_distance_cos(embedding, vector32(?))
		FROM vectors
		ORDER BY vector_distance_cos(embedding, vector32(?))
		ASC LIMIT 3;`, vectorStr, vectorStr)
	if err != nil {
		return nil, err
	}

	return similarItems, nil
}

func (s *DatabaseService) GetSettings() (*Settings, error) {
	var settings Settings
	err := s.db.Get(&settings, "SELECT url, llm, embedding_model FROM settings")
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (s *DatabaseService) UpdateSettings(url string, llm string, embedding string) error {
	_, err := s.db.Exec("UPDATE settings SET url=?, llm=?, embedding_model=? WHERE id=1", url, llm, embedding)
	if err != nil {
		return err
	}
	return nil
}
