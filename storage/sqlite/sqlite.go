package sqlite

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/victorfernandesraton/bushido/bushido"
)

func New() (bushido.LocalStorage, error) {
	db, err := sql.Open("sqlite3", "sqlite-bushido.db")
	if err != nil {
		return nil, err
	}

	st := &StorageSqlite{db: db}
	if err := st.CreateTables(); err != nil {
		return nil, err
	}
	return st, nil
}

type StorageSqlite struct {
	db *sql.DB
}

func (s *StorageSqlite) CreateTables() error {
	err := s.createTableContent()
	if err != nil {
		return err
	}

	err = s.createTableChapters()
	if err != nil {
		return err
	}

	err = s.createTablePages()
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageSqlite) createTableContent() error {
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

func (s *StorageSqlite) createTableChapters() error {
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
func (s *StorageSqlite) createTablePages() error {
	query := `CREATE TABLE IF NOT EXISTS page (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chapter_id INTEGER NOT NULL,
		content_id INTEGER NOT NULL,
		source TEXT NOT NULL,
		page_index INTEGER NOT NULL,
		link TEXT NOT NULL,
		FOREIGN KEY (content_id) REFERENCES content(id),
		FOREIGN KEY (chapter_id) REFERENCES chapter(id),
		UNIQUE( content_id, chapter_id, page_index, source)
	);`

	_, err := s.db.Exec(query)
	return err

}
func (s *StorageSqlite) Add(content bushido.Content) error {
	query := `INSERT INTO content (external_id, title, link, source, description, author) values (?, ? , ? , ?, ? , ?)`

	_, err := s.db.Exec(query, content.ExternalId, content.Title, content.Link, content.Source.ID, content.Description, content.Author)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageSqlite) FindById(id int) (*bushido.Content, error) {
	query := `select id, external_id, title, link, source, description, author from content where id = ?`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var contents []bushido.Content
	for rows.Next() {
		var content bushido.Content
		if err := rows.Scan(&content.ID, &content.ExternalId, &content.Title, &content.Link, &content.Source.ID, &content.Description, &content.Author); err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}

	if len(contents) != 1 {
		return nil, errors.New("content not found")
	}

	return &contents[0], nil
}

func (s *StorageSqlite) FindByLink(link string) (*bushido.Content, error) {
	query := `select id, external_id, title, link, source, description, author from content where link = ?`

	rows, err := s.db.Query(query, link)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var contents []bushido.Content
	for rows.Next() {
		var content bushido.Content
		if err := rows.Scan(&content.ID, &content.ExternalId, &content.Title, &content.Link, &content.Source.ID, &content.Description, &content.Author); err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}

	if len(contents) != 1 {
		return nil, errors.New("content not found")
	}

	return &contents[0], nil
}

func (s *StorageSqlite) ListByName(name string) ([]bushido.Content, error) {
	query := `select id, external_id, title, link, source, description, author from content where LOWER(title) LIKE '%' || LOWER(?) || '%'`

	rows, err := s.db.Query(query, name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var contents []bushido.Content
	for rows.Next() {
		var content bushido.Content
		if err := rows.Scan(&content.ID, &content.ExternalId, &content.Title, &content.Link, &content.Source.ID, &content.Description, &content.Author); err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}

	return contents, nil
}

func (s *StorageSqlite) AppendChapter(content bushido.Content, chapters []bushido.Chapter) error {
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
		_, err := stmt.Exec(c.ExternalId, c.Title, c.Link, c.Content.Source, content.ID)
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

func (s *StorageSqlite) ListChaptersByContentId(contentId int) ([]bushido.Chapter, error) {
	var result []bushido.Chapter
	query := `SELECT
		id,
		external_id,
		title,
		link,
		source,
		content_id
	FROM chapter WHERE content_id = ?`

	rows, err := s.db.Query(query, contentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var chapter bushido.Chapter
		var content bushido.Content
		if err := rows.Scan(&chapter.ID, &chapter.ExternalId, &chapter.Title, &chapter.Link, &content.Source.ID, &content.ExternalId); err != nil {
			return nil, err
		}
		chapter.Content = &content
		result = append(result, chapter)
	}

	return result, nil
}

func (s *StorageSqlite) AppendPages(chapter bushido.Chapter, pages []bushido.Page) error {
	query := `INSERT INTO page (
		content_id,
		chapter_id,
		source,
		page_index,
		link
	)
    VALUES (?, ?, ?, ?, ?)
    ON CONFLICT(content_id, chapter_id, page_index, source) DO UPDATE SET
        link = excluded.link`
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	for idx, p := range pages {
		_, err := stmt.Exec(chapter.Content.ID, chapter.ID, chapter.Content.Source, idx, p)
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

func (s *StorageSqlite) FindChapterById(id int) (*bushido.Chapter, error) {
	query := `SELECT
		id,
		external_id,
		title,
		link,
		source,
		content_id
	FROM chapter WHERE id = ?`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var chapters []bushido.Chapter
	for rows.Next() {
		var content bushido.Content
		var chapter bushido.Chapter
		if err := rows.Scan(&chapter.ID, &chapter.ExternalId, &chapter.Title, &chapter.Link, &content.Source.ID, &content.ID); err != nil {
			return nil, err
		}
		chapter.Content = &content
		chapters = append(chapters, chapter)
	}

	if len(chapters) != 1 {
		return nil, errors.New("chapter not found")
	}

	return &chapters[0], nil
}
