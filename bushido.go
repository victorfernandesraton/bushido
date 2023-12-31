package bushido

type Content struct {
	ID            int
	ExternalId    string
	Title         string
	Link          string
	Source        string
	Description   string
	Author        string
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

type Client interface {
	Search(query string) ([]Content, error)
	Chapters(link string, recursive bool) ([]Chapter, error)
	Pages(contentId string, chapterId string) ([]Page, error)
	Info(link string) (*Content, error)
	Source() string
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
