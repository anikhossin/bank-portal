package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func openDB(path string) *sql.DB {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	schema := `
	CREATE TABLE IF NOT EXISTS cases (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		case_id TEXT UNIQUE NOT NULL,
		iban TEXT NOT NULL,
		aes_reference TEXT NOT NULL,
		issue_type TEXT NOT NULL CHECK (issue_type IN ('failed_transfer', 'dispute')),
		issue_description TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'open',
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		token TEXT UNIQUE NOT NULL,
		token_used INTEGER NOT NULL DEFAULT 0
	);`
	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	// add new columns to existing DBs (ignore error if already present)
	db.Exec("ALTER TABLE cases ADD COLUMN verwendungszweck TEXT NOT NULL DEFAULT ''")
	db.Exec("ALTER TABLE cases ADD COLUMN betrag TEXT NOT NULL DEFAULT ''")
	return db
}
