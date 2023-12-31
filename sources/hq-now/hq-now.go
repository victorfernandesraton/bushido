package hqnow

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/victorfernandesraton/bushido/bushido"
)

type HqNow struct {
	*bushido.Source
	httpClient *http.Client
}
type GetHqsByNameResponse struct {
	Data struct {
		Items []GetHqsByName `json:"getHqsByName"`
	} `json:"data"`
}

type GetHqsByName struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	EditoraID        int    `json:"editoraId"`
	Status           string `json:"status"`
	PublisherName    string `json:"publisherName"`
	ImpressionsCount int    `json:"impressionsCount"`
}
type HqNowBody struct {
	OperationName string `json:"operationName"`
	Variables     any    `json:"variables"`
	Query         string `json:"query"`
}

type HqNowVariablesByName struct {
	Name string `json:"name"`
}

func NewClient(httpClient *http.Client) *HqNow {
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
	if &body.Data == nil {
		err = errors.New("Response not have data")
		return
	}
	result := &body.Data
	for _, item := range result.Items {
		content = append(content, *c.parseSearchResult(&item))
	}
	return
}

func (c *HqNow) parseSearchResult(item *GetHqsByName) (result *bushido.Content) {
	result = &bushido.Content{
		ExternalId: fmt.Sprintf("%v", item.ID),
		Title:      item.Name,
		Language:   "pt-BR",
		Source:     c.Source,
	}
	return
}

/*
	Chapters(link string, recursive bool) ([]Chapter, error)
	Pages(contentId string, chapterId string) ([]Page, error)
	Info(link string) (*Content, error)
	Source() string
*/
