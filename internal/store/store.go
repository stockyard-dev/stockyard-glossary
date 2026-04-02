package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct { db *sql.DB }

type Term struct {
	ID           string   `json:"id"`
	Term         string   `json:"term"`
	Definition   string   `json:"definition"`
	Category     string   `json:"category"`
	CreatedAt    string   `json:"created_at"`
}

func Open(dataDir string) (*DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	dsn := filepath.Join(dataDir, "glossary.db") + "?_journal_mode=WAL&_busy_timeout=5000"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS terms (
			id TEXT PRIMARY KEY,\n\t\t\tterm TEXT DEFAULT '',\n\t\t\tdefinition TEXT DEFAULT '',\n\t\t\tcategory TEXT DEFAULT '',
			created_at TEXT DEFAULT (datetime('now'))
		)`)
	if err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }

func (d *DB) Create(e *Term) error {
	e.ID = genID()
	e.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	_, err := d.db.Exec(`INSERT INTO terms (id, term, definition, category, created_at) VALUES (?, ?, ?, ?, ?)`,
		e.ID, e.Term, e.Definition, e.Category, e.CreatedAt)
	return err
}

func (d *DB) Get(id string) *Term {
	row := d.db.QueryRow(`SELECT id, term, definition, category, created_at FROM terms WHERE id=?`, id)
	var e Term
	if err := row.Scan(&e.ID, &e.Term, &e.Definition, &e.Category, &e.CreatedAt); err != nil {
		return nil
	}
	return &e
}

func (d *DB) List() []Term {
	rows, err := d.db.Query(`SELECT id, term, definition, category, created_at FROM terms ORDER BY created_at DESC`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []Term
	for rows.Next() {
		var e Term
		if err := rows.Scan(&e.ID, &e.Term, &e.Definition, &e.Category, &e.CreatedAt); err != nil {
			continue
		}
		result = append(result, e)
	}
	return result
}

func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM terms WHERE id=?`, id)
	return err
}

func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM terms`).Scan(&n)
	return n
}
