package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	"github.com/mirhijinam/url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB // connection to the database
}

func New(storagePath string) (*Storage, error) {
	fmt.Println(storagePath)
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"
	stmt, err := s.db.Prepare(`
		INSERT INTO url(url, alias) VALUES(?, ?)
	`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"
	stmt, err := s.db.Prepare(`
		SELECT url FROM url WHERE alias = ?
	`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var urlToGet string

	err = stmt.QueryRow(alias).Scan(&urlToGet)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return urlToGet, nil
}

// TODO: implement method
// func (s *Storage) DeleteURL(alias string) error {}
