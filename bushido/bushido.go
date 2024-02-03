package bushido

import (
	"net/url"
)

type Content struct {
	ID            string
	Title         string
	Link          string
	Source        *Source
	Description   string
	Author        string
	Language      string
	TotalChapters int
	Cover         *url.URL
}

type Chapter struct {
	ID        string
	Source    string
	Title     string
	Link      string
	ContentId string
}

type Page string

type Source struct {
	Domain    string
	ID        string
	Active    bool
	Languages []string
}

type Client interface {
	Search(query string) ([]Content, error)
	Chapters(id string) ([]Chapter, error)
	Pages(contentId string, chapterId string) ([]Page, error)
	Info(id string) (*Content, error)
}

type LocalStorage interface {
	Add(Content) error
	FindById(int) (*Content, error)
	FindByLink(string) (*Content, error)
	ListByName(string) ([]Content, error)
	AppendChapter(Content, []Chapter) error
	ListChaptersByContentId(int) ([]Chapter, error)
	FindChapterById(int) (*Chapter, error)
	AppendPages(Chapter, []Page) error
}
