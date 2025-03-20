-- name: InsertSaveGame :exec
INSERT INTO save_games (name) VALUES (?);

-- name: InsertChatMessage :exec
INSERT INTO chat_messages (role, content, save_game_id) VALUES (?, ?, ?);

-- name: GetSaveGames :many
SELECT * FROM save_games;

-- name: GetSaveGame :one
SELECT * FROM save_games WHERE id = ?;

-- name: GetChatMessages :many
SELECT * FROM chat_messages WHERE save_game_id = ?;

-- name: DeleteSaveGame :exec
DELETE FROM save_games WHERE id = ?;

