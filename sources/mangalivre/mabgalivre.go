package mangalivre

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/victorfernandesraton/bushido"
)

type MangaLivre struct {
}

type seriesItem struct {
	IDSerie        int    `json:"id_serie,omitempty"`
	Name           string `json:"name,omitempty"`
	Label          string `json:"label,omitempty"`
	Score          string `json:"score,omitempty"`
	Value          string `json:"value,omitempty"`
	Author         string `json:"author,omitempty"`
	Artist         string `json:"artist,omitempty"`
	Cover          string `json:"cover,omitempty"`
	CoverThumb     string `json:"cover_thumb,omitempty"`
	CoverAvif      string `json:"cover_avif,omitempty"`
	CoverThumbAvif string `json:"cover_thumb_avif,omitempty"`
	Link           string `json:"link,omitempty"`
	IsComplete     bool   `json:"is_complete,omitempty"`
}

type seriesResponse struct {
	Series *[]seriesItem `json:"series"`
}

type chapterItem struct {
	IDSerie     int    `json:"id_serie,omitempty"`
	IDChapter   int    `json:"id_chapter,omitempty"`
	Name        string `json:"name,omitempty"`
	ChapterName string `json:"chapter_name,omitempty"`
	Number      string `json:"number,omitempty"`
}

type chapterResponse struct {
	Chapters *[]chapterItem `json:"chapters"`
}

type page struct {
	Legacy string `json:"legacy,omitempty"`
	Avif   string `json:"avif,omitempty"`
}

type pageResponse struct {
	Images *[]page
}

func (source *MangaLivre) Search(query string) (*[]bushido.BasicContent, error) {

	formData := url.Values{}
	formData.Set("search", query)

	req, err := http.NewRequest("POST", "https://mangalivre.net/lib/search/series.json", bytes.NewReader([]byte(formData.Encode())))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request status error, expect 200, got %v", res.StatusCode)
	}

	if err != nil {
		return nil, err
	}

	var data seriesResponse
	var result []bushido.BasicContent

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		if strings.Contains(err.Error(), "json: cannot unmarshal bool into") {
			return &result, nil
		}
		return nil, err
	}

	if data.Series == nil {
		return &result, nil
	}

	for _, v := range *data.Series {
		result = append(result, bushido.BasicContent{
			ExternalId: fmt.Sprintf("%d", v.IDSerie),
			Title:      v.Name,
			Source:     "mangalivre",
			Link:       v.Link,
		})
	}

	return &result, nil

}

func (source *MangaLivre) parseUrlToId(url string) (int, error) {
	paths := strings.Split(url, "/")
	if len(paths) < 6 {
		return 0, fmt.Errorf("not valid url")
	}

	value, err := strconv.Atoi(paths[5])
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (source *MangaLivre) Chapters(link string, page int) (*[]bushido.Chapter, error) {
	id, err := source.parseUrlToId(link)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://mangalivre.net/series/chapters_list.json?page=%v&id_serie=%v", page, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request status error, expect 200, got %v", res.StatusCode)
	}

	if err != nil {
		return nil, err
	}

	var data chapterResponse
	var result []bushido.Chapter

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		if strings.Contains(err.Error(), "json: cannot unmarshal bool into") {
			return &result, nil
		}
		return nil, err
	}

	if data.Chapters == nil {
		return &result, nil
	}

	for _, v := range *data.Chapters {
		result = append(result, bushido.Chapter{
			ExternalId: fmt.Sprintf("%d", v.IDSerie),
			Title:      fmt.Sprintf("%s - %s", v.Name, v.Number),
		})
	}

	return &result, nil
}
func (source *MangaLivre) Pages(contentId string, chapterId string) (*[]bushido.Page, error) {

	url := fmt.Sprintf("https://mangalivre.net/leitor/pages/%s.jso", chapterId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request status error, expect 200, got %v", res.StatusCode)
	}

	if err != nil {
		return nil, err
	}

	var data pageResponse
	var result []bushido.Page

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		if strings.Contains(err.Error(), "json: cannot unmarshal bool into") {
			return &result, nil
		}
		return nil, err
	}

	if data.Images == nil {
		return &result, nil
	}

	for _, v := range *data.Images {
		result = append(result, bushido.Page(v.Legacy))
	}

	return &result, nil
}

func (source *MangaLivre) Info(link string) (*bushido.Content, error) {

	if !strings.Contains(link, "https://mangalivre.net/") {
		return nil, fmt.Errorf("not valid url")
	}
	doc, err := htmlquery.LoadURL(link)
	if err != nil {
		return nil, err
	}

	descriptionNosw := htmlquery.FindOne(doc, "//html/body/div[5]/div/div[3]/div[5]/div[2]/span[3]/span")
	titleNode := htmlquery.FindOne(doc, "//html/body/div[5]/div/div[3]/div[5]/div[2]/span[1]/h1")
	totalChaptersStr := htmlquery.FindOne(doc, "//html/body/div[5]/div/div[4]/div[3]/h2/span")
	fmt.Println(totalChaptersStr)
	totalChapters, err := strconv.Atoi(totalChaptersStr.FirstChild.Data)
	if err != nil {
		return nil, errors.New("not valid find html total chapters")
	}

	return &bushido.Content{
		BasicContent: bushido.BasicContent{
			Title:  titleNode.FirstChild.Data,
			Link:   link,
			Source: "mangalivre",
		},
		TotalChapters: int64(totalChapters),
		Description:   descriptionNosw.FirstChild.Data,
	}, nil
}

// Install(link string) error
// Sync() error
// List(query string) (error []Content)
// Remove(id uint64) error
