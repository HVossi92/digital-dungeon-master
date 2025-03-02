package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/tursodatabase/go-libsql"
)

// VectorService represents the service responsible for the vector database.
type VectorService struct {
	db *sql.DB
}

type VectorItem struct {
	Text      string
	Embedding []byte
}

type Settings struct {
	URL       string
	LLM       string
	Embedding string
}

// SetUDatabaseService creates and initializes a new VectorDBService.
func SetUDatabaseService(dbPath string, overwrite bool) (*VectorService, error) {
	if overwrite {
		log.Println("Overwriting existing database (if it exists)")
		if err := os.Remove(dbPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("failed to remove existing database: %w", err)
		}
	}

	// Create VectorService instance first
	vectorService := &VectorService{db: nil}
	db, err := vectorService.createDb(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	vectorService.db = db

	// Then call ensureVectorTableExists on the instance
	if err := vectorService.createVectorTable(); err != nil {
		db.Close() // Close the connection if table creation fails
		return nil, fmt.Errorf("failed to ensure vector table exists: %w", err)
	}

	if err := vectorService.createSettingsTable(); err != nil {
		db.Close() // Close the connection if table creation fails
		return nil, fmt.Errorf("failed to ensure settings table exists: %w", err)
	}

	return vectorService, nil
}

// Close closes the database connection.  Good practice to add a Close method.
func (s *VectorService) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// GetDB returns the underlying sql.DB connection (for use within the service package).
func (s *VectorService) GetDB() *sql.DB {
	return s.db
}

func (s *VectorService) createDb(dbPath string) (*sql.DB, error) {
	// Connect to embedded libSQL
	db, err := sql.Open("libsql", "file:"+dbPath)
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
func (s *VectorService) createVectorTable() error {
	_, err := s.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS vectors (id INTEGER PRIMARY KEY, title TEXT, text TEXT, embedding F32_BLOB(%d))", 768))
	if err != nil {
		return err
	}

	return nil
}

func (s *VectorService) createSettingsTable() error {
	// Create table if not exists
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS settings (
		id INTEGER PRIMARY KEY, 
		url TEXT, 
		llm TEXT, 
		embedding_model TEXT
	) STRICT`)
	if err != nil {
		return fmt.Errorf("failed to create settings table: %w", err)
	}

	// Check if table is empty
	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM settings").Scan(&count)
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

// StoreChunkAndEmbedding saves a text chunk and its embedding to the SQLite vector database.
func (s *VectorService) StoreChunkAndEmbedding(chunk string, embedding []float32) error {
	if s.db == nil {
		return fmt.Errorf("database connection is nil")
	}

	title := chunk
	if len(chunk) > 8 {
		title = chunk[:8]
	}

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
		`INSERT INTO vectors (title, text, embedding) 
         VALUES (?, ?, vector32(?))`,
		title,
		chunk,
		vectorStr,
	)
	if err != nil {
		return err
	}
	return nil
}

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
func (s *VectorService) ChunkText(text string, chunkSize int, chunkOverlap int) ([]string, error) {
	if chunkSize <= 0 {
		return nil, fmt.Errorf("chunkSize must be greater than 0")
	}
	if chunkOverlap < 0 {
		return nil, fmt.Errorf("chunkOverlap must be non-negative")
	}
	if chunkOverlap >= chunkSize {
		return nil, fmt.Errorf("chunkOverlap must be less than chunkSize")
	}

	sentences := splitIntoSentences(text) // Split text into sentences
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

// splitIntoSentences splits a text into sentences using a simple regex.
// Note: This is a basic sentence splitter and might not be perfect for all cases.
// For more robust sentence splitting, consider using NLP libraries.
func splitIntoSentences(text string) []string {
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

// ReadAllVectors reads all data from the vector DB table.
func (s *VectorService) ReadAllVectors() (string, error) {
	if s.db == nil {
		return "", fmt.Errorf("database connection is nil")
	}

	var builder strings.Builder
	const (
		maxTextLength  = 40
		maxEmbedLength = 40
	)

	rows, err := s.db.Query(`
		SELECT id, title, text, vector_extract(embedding) 
		FROM vectors
		ORDER BY id ASC`)
	if err != nil {
		return "", fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var (
		id        int
		title     string
		text      string
		embedding string
	)

	for rows.Next() {
		err := rows.Scan(&id, &title, &text, &embedding)
		if err != nil {
			return builder.String(), fmt.Errorf("scan failed: %w", err)
		}

		// Truncate fields
		displayText := text
		if len(text) > maxTextLength {
			displayText = text[:maxTextLength-3] + "..."
		}

		builder.WriteString(fmt.Sprintf("%-4d - %-40s | ", id, displayText))
	}

	if err := rows.Err(); err != nil {
		return builder.String(), fmt.Errorf("row iteration error: %w", err)
	}

	return builder.String(), nil
}

// FindSimilarVectors queries the vector DB for vectors similar to the given embedding.
func (s *VectorService) FindSimilarVectors(queryEmbedding []float32) ([]VectorItem, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection is nil in VectorService")
	}

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

	rows, err := s.db.Query(
		`SELECT title, text, vector_extract(embedding),
       vector_distance_cos(embedding, vector32(?))
		FROM vectors
		ORDER BY
       vector_distance_cos(embedding, vector32(?))
		ASC LIMIT 3;`, vectorStr, vectorStr)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Iterate through results
	var (
		title     string
		text      string
		embedding string
		distance  float64
	)

	var similarItems []VectorItem
	for rows.Next() {
		err := rows.Scan(&title, &text, &embedding, &distance)
		if err != nil {
			return nil, err
		}

		// Format the output
		fmt.Printf("%-20s | %.4f\n",
			text,
			distance)
		item := VectorItem{
			Text:      text,
			Embedding: []byte(embedding),
		}
		similarItems = append(similarItems, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return similarItems, nil
}

func (s *VectorService) GetSettings() (*Settings, error) {
	var settings Settings
	err := s.db.QueryRow("SELECT url, llm, embedding_model FROM settings").Scan(&settings.URL, &settings.LLM, &settings.Embedding)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (s *VectorService) UpdateSettings(url string, llm string, embedding string) error {
	_, err := s.db.Exec("UPDATE settings SET url=?, llm=?, embedding_model=? WHERE id=1", url, llm, embedding)
	if err != nil {
		return err
	}
	return nil
}
