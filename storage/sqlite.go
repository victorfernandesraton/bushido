package storage

import (
	"database/sql"
	"errors"
	"log"

	"github.com/victorfernandesraton/bushido"
)

type StorageSqlte struct {
	db *sql.DB
}

func New(db *sql.DB) *StorageSqlte {
	return &StorageSqlte{db: db}
}

func (s *StorageSqlte) CreateTables() error {
	err := s.createTableContent()
	if err != nil {
		return err
	}

	err = s.createTableChapters()
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageSqlte) createTableContent() error {
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

func (s *StorageSqlte) createTableChapters() error {
	query := `
    CREATE TABLE IF NOT EXISTS chapter (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		external_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		link TEXT NOT NULL,
		source TEXT NOT NULL,
		content_id INTEGER NOT NULL,
		FOREIGN KEY (content_id) REFERENCES content(id),
		UNIQUE(external_id, content_id)
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

func (s *StorageSqlte) ListByName(name string) ([]bushido.Content, error) {
	query := `select external_id, title, link, source, description, author from content where LOWER(title) LIKE '%' || LOWER(?) || '%'`

	rows, err := s.db.Query(query, name)
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

	return contents, nil
}

func (s *StorageSqlte) AppendChapter(id int, chapters []bushido.Chapter) error {
	query := `INSERT INTO chapter (external_id, title, link, source, content_id)
    VALUES (?, ?, ?, ?, ?)
    ON CONFLICT(external_id, content_id) DO UPDATE SET
        title = excluded.title,
        link = excluded.link,
        source = excluded.source,
        content_id = excluded.content_id`
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	for _, c := range chapters {
		log.Println("append chapter ", c.Title, c.ExternalId)
		_, err := stmt.Exec(c.ExternalId, c.Title, c.Link, c.Content.Source, id)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil

}

func (s *StorageSqlte) ListChaptersByContentId(contentId int) ([]bushido.Chapter, error) {
	var result []bushido.Chapter
	rows, err := s.db.Query(`
        SELECT
            external_id,
            title,
            link,
            source,
            content_id
        FROM chapter WHERE content_id = ?`,
		contentId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var chapter bushido.Chapter
		var content bushido.Content
		if err := rows.Scan(&chapter.ExternalId, &chapter.Title, &chapter.Link, &content.Source, &content.ExternalId); err != nil {
			return nil, err
		}
		chapter.Content = &content
		result = append(result, chapter)
	}

	return result, nil
}
