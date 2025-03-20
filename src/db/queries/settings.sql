-- name: InsertSettings :exec
INSERT INTO settings (url, llm)
VALUES (?, ?);

-- name: GetSettings :one
SELECT * FROM settings;

-- name: UpdateSettings :exec
UPDATE settings SET url=?, llm=? WHERE id=1;