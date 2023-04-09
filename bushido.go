package bushido

type Content struct {
	ExternalId  string
	Title       string
	Link        string
	Description string
	Source      string
	Author      string
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
	Install(link string) error
	Sync() error
	List(query string) (error []Content)
	Remove(id uint64) error
}
