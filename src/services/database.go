package services

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDb(dbPath string, overwrite bool) *sql.DB {
	if overwrite {
		log.Println("Overwriting existing database (if it exists)")
		if err := os.Remove(dbPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}
	}

	// Connect to embedded libSQL
	connStr := "./database.db?_foreign_keys=on&_journal=WAL&_synchronous=NORMAL&_busy_timeout=5000&_cache_size=-64000"
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Connection failed:", err)
	}

	log.Println("Connected to local sqlite database")
	return db
}

func CloseDb(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// func (s *DatabaseService) createSettingsTable() error {
// 	// Create table if not exists
// 	s.db.MustExec(`CREATE TABLE IF NOT EXISTS settings (
// 		id INTEGER PRIMARY KEY,
// 		url TEXT,
// 		llm TEXT,
// 		embedding_model TEXT
// 	) STRICT`)

// 	// Check if table is empty
// 	var count int
// 	err := s.db.Get(&count, "SELECT COUNT(*) FROM settings")
// 	if err != nil {
// 		return fmt.Errorf("failed to check settings table: %w", err)
// 	}

// 	// Insert default values only if table is empty
// 	if count == 0 {
// 		_, err = s.db.Exec(`INSERT INTO settings (url, llm, embedding_model)
// 			VALUES (?, ?, ?)`,
// 			"http://192.168.178.105:11434",
// 			"gemma3:4b",
// 			"nomic-embed-text:latest")
// 		if err != nil {
// 			return fmt.Errorf("failed to insert default settings: %w", err)
// 		}
// 	}

// 	return nil
// }

// func (s *DatabaseService) createChatMessageTable() error {
// 	// Create table if not exists
// 	s.db.MustExec(`CREATE TABLE IF NOT EXISTS chat_messages (
// 		id INTEGER PRIMARY KEY,
// 		role TEXT,
// 		content TEXT,
// 		created TEXT DEFAULT CURRENT_TIMESTAMP
// 	) STRICT`)

// 	return nil
// }

// func (s *DatabaseService) createSaveGamesTables() error {
// 	// Create table if not exists
// 	s.db.MustExec(`CREATE TABLE IF NOT EXISTS save_games (
// 		id INTEGER PRIMARY KEY,
// 		name TEXT NOT NULL,
// 		created TEXT DEFAULT CURRENT_TIMESTAMP
// 	) STRICT`)

// 	s.db.MustExec(`CREATE TABLE IF NOT EXISTS chat_messages (
// 		id INTEGER PRIMARY KEY,
// 		role TEXT,
// 		content TEXT,
// 		save_game_id INTEGER,
// 		created TEXT DEFAULT CURRENT_TIMESTAMP,
// 		FOREIGN KEY (save_game_id) REFERENCES save_games(id) ON DELETE CASCADE
// 	) STRICT`)

// 	return nil
// }

// func (s *DatabaseService) GetSettings() (*Settings, error) {
// 	var settings Settings
// 	err := s.db.Get(&settings, "SELECT url, llm, embedding_model FROM settings")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &settings, nil
// }

// func (s *DatabaseService) UpdateSettings(url string, llm string, embedding string) error {
// 	_, err := s.db.Exec("UPDATE settings SET url=?, llm=?, embedding_model=? WHERE id=1", url, llm, embedding)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (s *DatabaseService) SaveGame(name string, messages []ChatMessage) error {
// 	result, err := s.db.Exec(
// 		`INSERT INTO save_games (name) VALUES (?)`,
// 		name,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return err
// 	}
// 	for _, msg := range messages {
// 		_, err := s.db.Exec(
// 			`INSERT INTO chat_messages (role, content, save_game_id) VALUES (?, ?, ?)`,
// 			msg.Role, msg.Content, id,
// 		)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (s *DatabaseService) GetSaveGames() ([]SaveGame, error) {
// 	var saves []SaveGame
// 	err := s.db.Select(&saves, `
// 		SELECT id, name, created
// 		FROM save_games
// 		ORDER BY created DESC
// 	`)
// 	for i := range saves {
// 		date, err := helpers.GetDayMonthYearFrom(saves[i].Created)
// 		if err != nil {
// 			return nil, err
// 		}
// 		saves[i].Created = date
// 	}
// 	return saves, err
// }

// func (s *DatabaseService) DeleteSaveGame(id int) error {
// 	_, err := s.db.Exec(`DELETE FROM save_games WHERE id = ?`, id)
// 	return err
// }

// func (s *DatabaseService) LoadSaveGame(id int) ([]ChatMessage, error) {
// 	var save SaveGame
// 	fmt.Println("LoadSaveGame")
// 	err := s.db.Get(&save, `
// 		SELECT id, name, created
// 		FROM save_games
// 		WHERE id = ?
// 	`, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var messages []ChatMessage
// 	err = s.db.Select(&messages, `
// 		SELECT role, content
// 		FROM chat_messages
// 		WHERE save_game_id = ?
// 		ORDER BY created ASC
// 	`, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return messages, err
// }
