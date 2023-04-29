package storage

import (
	"database/sql"
	"errors"

	"github.com/victorfernandesraton/bushido"
)

type StorageSqlte struct {
	db *sql.DB
}

func New(db *sql.DB) *StorageSqlte {
	return &StorageSqlte{db: db}
}

func (s *StorageSqlte) CreateTables() error {
	query := `
    CREATE TABLE IF NOT EXISTS content(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		external_id INTEGER NOT NULL,
		title TEXT NOT NULL,
        link TEXT NOT NULL,
		source TEXT,
		description TEXT,
		author TEXT
    );`

	_, err := s.db.Exec(query)
	return err

}

func (s *StorageSqlte) Add(content bushido.Content) error {
	query := `INSERT INTO content (external_id, title, link, source, description, author) values (?, ? , ? , ?, ? , ?)`

	_, err := s.db.Exec(query, content.ExternalId, content.Title, content.Link, content.Source, content.Description, content.Author)
	if err != nil {

		return err
	}
	return nil
}

func (s *StorageSqlte) FindById(id int) (*bushido.Content, error) {
	query := `select external_id, title, link, source, description, author from content where id = ?`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var contents []bushido.Content
	for rows.Next() {
		var content bushido.Content
		if err := rows.Scan(&content.ExternalId, &content.Title, &content.Link, &content.Source, &content.Description, &content.Author); err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}

	if len(contents) != 1 {
		return nil, errors.New("content not found")
	}

	return &contents[0], nil
}

func (s *StorageSqlte) FindByLink(link string) (*bushido.Content, error) {
	query := `select external_id, title, link, source, description, author from content where link = ?`

	rows, err := s.db.Query(query, link)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var contents []bushido.Content
	for rows.Next() {
		var content bushido.Content
		if err := rows.Scan(&content.ExternalId, &content.Title, &content.Link, &content.Source, &content.Description, &content.Author); err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}

	if len(contents) != 1 {
		return nil, errors.New("content not found")
	}

	return &contents[0], nil
}
