package hqnow_test

import (
	"net/http"
	"testing"

	"github.com/victorfernandesraton/bushido/bushido/helpers/assertion"
	hqnow "github.com/victorfernandesraton/bushido/sources/hq-now"
)

func TestSearchIntegation(t *testing.T) {
	t.Run("search by spider", func(t *testing.T) {
		client := hqnow.NewClient(http.DefaultClient)

		result, err := client.Search("spider")
		assertion.AssertEqual(t, nil, err)
		assertion.AssertEqual(t, 4, len(result))
		t.Log(result)
	})

	t.Run("search by dibrish world", func(t *testing.T) {
		client := hqnow.NewClient(http.DefaultClient)

		result, err := client.Search("sskadjj")
		assertion.AssertEqual(t, nil, err)
		assertion.AssertEqual(t, 0, len(result))
		t.Log(result)
	})
}
