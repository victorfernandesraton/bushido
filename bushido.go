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
	TotalChapters int64
}

type Chapter struct {
	ExternalId string
	Title      string
	Content    *Content
}

type Page string

type Client interface {
	Search(query string) (*[]Content, error)
	Chapters(link string) (*[]Chapter, error)
	Pages(contentId string, chapterId string) (*[]Page, error)
	Info(link string) (*Content, error)
	// Deepends of sqlite
	Install(link string) error
	Sync(link string) error
	List() (error []Content)
	Remove(id uint64) error
}
