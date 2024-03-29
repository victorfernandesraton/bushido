package bushido

type Content struct {
	ID            int
	ExternalId    string
	Title         string
	Link          string
	Source        *Source
	Description   string
	Author        string
	Language      string
	TotalChapters int
}

type Chapter struct {
	ID         int
	ExternalId string
	Title      string
	Link       string
	Content    *Content
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
	Chapters(link string, recursive bool) ([]Chapter, error)
	Pages(contentId string, chapterId string) ([]Page, error)
	Info(link string) (*Content, error)
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
