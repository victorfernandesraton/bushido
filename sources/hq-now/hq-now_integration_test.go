package hqnow_test

import (
	"net/http"
	"testing"

	"github.com/victorfernandesraton/bushido/helpers/assertion"
	hqnow "github.com/victorfernandesraton/bushido/sources/hq-now"
)

var client = hqnow.NewClient(http.DefaultClient)

func TestSearchIntegation(t *testing.T) {
	t.Run("search by spider", func(t *testing.T) {
		result, err := client.Search("spider")
		assertion.AssertEqual(t, nil, err)
		assertion.AssertEqual(t, 4, len(result))
		t.Log(result)
	})

	t.Run("search by dibrish world", func(t *testing.T) {
		result, err := client.Search("sskadjj")
		assertion.AssertEqual(t, nil, err)
		assertion.AssertEqual(t, 0, len(result))
		t.Log(result)
	})
}

func TestInfoIntegration(t *testing.T) {

	t.Run("get by know id", func(t *testing.T) {
		result, err := client.Info("602")
		assertion.AssertEqual(t, nil, err)
		assertion.AssertNotNil(t, result)
		assertion.AssertEqual(t, "Amazing Spider Man", result.Title)
		t.Log(result)
	})

	t.Run("not found content by id", func(t *testing.T) {
		t.Skip("error in nill comparation")
		result, err := client.Info("1111111")
		assertion.AssertEqual(t, nil, err)
		assertion.AssertDeepEqual(t, result, nil)
		t.Log(result)
	})
}

func TestChapterIntegration(t *testing.T) {

	t.Run("get chapters by hq id", func(t *testing.T) {
		result, err := client.Chapters("602")
		assertion.AssertEqual(t, nil, err)
		assertion.AssertNotNil(t, result)
		assertion.AssertEqual(t, "10632", result[1].ID)
		assertion.AssertEqual(t, "602", result[1].ContentId)
		assertion.AssertEqual(t, "Anual 1", result[1].Title)
		t.Log(result[1])
	})

	t.Run("not found hqid", func(t *testing.T) {
		result, err := client.Chapters("1111111")
		assertion.AssertEqual(t, nil, err)
		assertion.AssertEqual(t, 0, len(result))
		t.Log(result)
	})
}
