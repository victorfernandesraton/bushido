package bushido

type BasicContent struct {
	ExternalId string
	Title      string
	Link       string
	Source     string
}

type Content struct {
	BasicContent
	Description   string
	Author        string
	TotalChapters int
}

type Chapter struct {
	ExternalId string
	Title      string
	Link       string
	Content    *Content
}

type Page string

type Client interface {
	Search(query string) ([]Content, error)
	Chapters(link string, recursive bool) ([]Chapter, error)
	Pages(contentId string, chapterId string) ([]Page, error)
	Info(link string) (*Content, error)
	Source() string
	// Deepends of sqlite
	// Install(link string) error
	// Sync(link string) error
	// List() ([]Content, error)
	// Remove(id uint64) error
}
