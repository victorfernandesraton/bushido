package hqnow

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/victorfernandesraton/bushido/bushido"
)

type HqNow struct {
	*bushido.Source
	httpClient *http.Client
}
type HqNowResponse struct {
	Data any `json:"data"`
}

type Chapter struct {
	Name      string `json:"name"`
	ID        int    `json:"id"`
	HqId      int    `json:"hqId"`
	Number    string `json:"number"`
	UpdatedAt *time.Time
}
type HqInfo struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	EditoraID     int    `json:"editoraId"`
	Synopsis      string `json:"synopsis"`
	Status        string `json:"status"`
	PublisherName string `json:"publisherName"`
	HqCover       string `json:"hqCover"`
}
type GetHqsByNameResponse struct {
	Data struct {
		Items []HqInfo `json:"getHqsByName"`
	} `json:"data"`
}

type GetHqsByIdResponse struct {
	Data struct {
		Items []GetHqsByID `json:"getHqsById"`
	} `json:"data"`
}

type GetHqChapters struct {
	Data struct {
		Items []Chapter `json:"getChaptersByHqId"`
	}
}

type GetHqsByID struct {
	HqInfo
	Chapters []Chapter `json:"capitulos"`
}
type HqNowBody struct {
	OperationName string `json:"operationName"`
	Variables     any    `json:"variables"`
	Query         string `json:"query"`
}

type HqNowVariablesByName struct {
	Name string `json:"name"`
}

type HqNowVariablesById struct {
	ID int64 `json:"id"`
}

func NewClient(httpClient *http.Client) bushido.Client {
	return &HqNow{
		httpClient: httpClient,
		Source: &bushido.Source{
			Active:    true,
			Domain:    "hq-now.com",
			ID:        "hqnow",
			Languages: []string{"pt-BR"},
		},
	}
}

func (c *HqNow) Search(query string) (content []bushido.Content, err error) {
	data := &HqNowBody{
		OperationName: "getHqsByName",
		Query:         `query getHqsByName($name: String!) {  getHqsByName(name: $name) {    id    name    editoraId    status    publisherName    impressionsCount  }}`,
		Variables: &HqNowVariablesByName{
			Name: query,
		},
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return
	}

	res, err := c.httpClient.Post("https://admin.hq-now.com/graphql", "application/json", bytes.NewBuffer(payload))

	if res.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Error in serach by hq in hqNow, status: %v", res.StatusCode))
	}
	body := &GetHqsByNameResponse{}
	json.NewDecoder(res.Body).Decode(&body)

	defer res.Body.Close()
	if &body.Data == nil {
		err = errors.New("Response not have data")
		return
	}

	result := &body.Data
	for _, item := range result.Items {
		updateItem, err := c.parseSearchResult(&item)
		if err != nil {
			err = errors.Join(errors.New("Error during parse data for conteent"), err)
		}
		content = append(content, *updateItem)
	}
	return
}

func (c *HqNow) parseSearchResult(item *HqInfo) (result *bushido.Content, err error) {
	cover, err := url.Parse(item.HqCover)
	result = &bushido.Content{
		ID:          fmt.Sprintf("%v", item.ID),
		Description: item.Synopsis,
		Title:       item.Name,
		Language:    "pt-BR",
		Source:      c.Source,
		Cover:       cover,
	}
	return
}

func (c *HqNow) Info(link string) (content *bushido.Content, err error) {
	id, err := strconv.ParseInt(link, 10, 64)
	data := &HqNowBody{
		OperationName: "getHqsById",
		Query:         `query getHqsById($id: Int!) {  getHqsById(id: $id) {    id    name    synopsis    editoraId    status    publisherName    hqCover    impressionsCount    capitulos {      name      id      number    }  }}`,
		Variables: &HqNowVariablesById{
			ID: id,
		},
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return
	}

	res, err := c.httpClient.Post("https://admin.hq-now.com/graphql", "application/json", bytes.NewBuffer(payload))

	if res.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Error in serach by hq in hqNow, status: %v", res.StatusCode))
	}
	body := &GetHqsByIdResponse{}
	json.NewDecoder(res.Body).Decode(&body)

	defer res.Body.Close()
	if &body.Data == nil {
		err = errors.New("Response not have data")
		return
	}
	result := &body.Data
	for _, item := range result.Items {
		content, err = c.parseSearchResult(&item.HqInfo)
		if err != nil {
			err = errors.Join(errors.New("Error during parse data for conteent"), err)
		}
		content.TotalChapters = len(item.Chapters)

	}
	return
}

func (c *HqNow) Chapters(link string) (chapters []bushido.Chapter, err error) {
	id, err := strconv.ParseInt(link, 10, 64)
	data := &HqNowBody{
		OperationName: "getChaptersByHqID",
		Query:         `query getChaptersByHqID($id: Int!) {getChaptersByHqId(hqId: $id){id name number updatedAt hqId}}`,
		Variables: &struct {
			ID int64 `json:"id"`
		}{
			ID: id,
		},
	}
	payload, err := json.Marshal(data)

	if err != nil {
		return
	}

	res, err := c.httpClient.Post("https://admin.hq-now.com/graphql", "application/json", bytes.NewBuffer(payload))

	if res.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Error in get chapters by hq in hqNow, status: %v", res.StatusCode))
	}

	body := &GetHqChapters{}
	json.NewDecoder(res.Body).Decode(&body)
	defer res.Body.Close()
	if &body.Data == nil {
		err = errors.New("Response not have data")
		return
	}

	result := &body.Data

	for _, item := range result.Items {
		chapter := c.parseChapter(item)
		chapters = append(chapters, chapter)

	}
	return
}

func (c *HqNow) parseChapter(item Chapter) bushido.Chapter {
	return bushido.Chapter{
		ID:        fmt.Sprintf("%v", item.ID),
		Title:     strings.TrimSpace(item.Name),
		ContentId: fmt.Sprintf("%v", item.HqId),
	}
}

func (c *HqNow) Pages(contentId string, chapterId string) (pages []bushido.Page, err error) {
	// TDD: wait impl
	return
}
