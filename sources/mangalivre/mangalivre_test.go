package mangalivre_test

import (
	"fmt"
	"testing"

	"github.com/victorfernandesraton/bushido/sources/mangalivre"
)

func Test_Search(t *testing.T) {
	stub := mangalivre.MangaLivre{}
	type args struct {
		query    string
		err      error
		hasItems bool
	}

	testcases := []struct {
		name string
		args args
	}{
		{
			name: "sucess naruto query",
			args: args{
				query: "naruto", err: nil,
				hasItems: true,
			},
		},
		{
			name: "error stranger query",
			args: args{
				query: "ngfjksdngkjlndskgjns", err: nil,
				hasItems: false,
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			data, err := stub.Search(tt.args.query)
			if err != tt.args.err {
				t.Errorf("expected err is %v, got %v", tt.args.err, err)
			}
			if data == nil && len(*data) == 0 {
				t.Errorf("expected data is not empty, got %v", data)
			}
		})
	}

}

func Test_Chapters(t *testing.T) {
	stub := mangalivre.MangaLivre{}
	type args struct {
		link         string
		err          error
		totalChapter int
		hasData      bool
	}

	testcases := []struct {
		name string
		args args
	}{
		{
			"success with solo leveling",
			args{
				link:         "https://mangalivre.net/manga/solo-leveling/7702",
				err:          nil,
				hasData:      true,
				totalChapter: 30,
			},
		},
		{
			"error with invalid url",
			args{
				link:         "https://www.google.com",
				err:          fmt.Errorf("not valid url"),
				totalChapter: 0,
				hasData:      false,
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			data, err := stub.Chapters(tt.args.link)

			if !tt.args.hasData {
				if err.Error() != tt.args.err.Error() {
					t.Errorf("expected err is %v, got %v", tt.args.err, err)
				}
				return
			}
			if err != nil {
				t.Errorf("expected err is nil, got %v", err)

			}
			if data == nil {
				t.Errorf("expected data is not empty, got %v", data)
			}
			if len(*data) != tt.args.totalChapter {
				t.Errorf("expected data is not empty, got %v", data)
			}
		})
	}
}

func Test_Pages(t *testing.T) {
	stub := mangalivre.MangaLivre{}
	type args struct {
		chapterId string
		// TODO: not using
		contnetId  string
		err        error
		hasItems   bool
		totalPages int
	}

	testcases := []struct {
		name string
		args args
	}{
		{
			"success in solo leveling cap 1",
			args{
				contnetId:  "7702",
				chapterId:  "455324",
				err:        nil,
				hasItems:   true,
				totalPages: 22,
			},
		},
		{
			"error not found chapter",
			args{
				contnetId:  "7702",
				chapterId:  "888888",
				err:        nil,
				hasItems:   true,
				totalPages: 0,
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			data, err := stub.Pages(tt.args.contnetId, tt.args.chapterId)
			if !tt.args.hasItems {
				if err.Error() != tt.args.err.Error() {
					t.Errorf("expected err is %v, got %v", tt.args.err, err)
				}
				return
			}
			if data == nil {
				t.Errorf("expected data is not empty, got %v", data)
			}
			if len(*data) != tt.args.totalPages {
				t.Errorf("expected data size is %d, got %d", tt.args.totalPages, len(*data))
			}
		})
	}

}

func Test_Info(t *testing.T) {
	stub := mangalivre.MangaLivre{}
	type args struct {
		url         string
		title       string
		description string
		err         error
		hasContent  bool
	}

	testcases := []struct {
		name string
		args args
	}{
		{
			"success with solo leveling search",
			args{
				url:         "https://mangalivre.net/manga/solo-leveling/7702",
				title:       "Solo Leveling",
				description: `Dez anos atrás, depois do "Portal" que conecta o mundo real com um mundo de montros se abriu, algumas pessoas comuns receberam o poder de caçar os monstros do portal. Eles são conhecidos como caçadores. Porém, nem todos os caçadores são fortes. Meu nome é Sung Jin-Woo, um caçador de rank E. Eu sou alguém que tem que arriscar a própria vida nas dungeons mais fracas, "O mais fraco do mundo". Sem ter nenhuma habilidade à disposição, eu mal consigo dinheiro nas dungeons de baixo nível... Ao menos até eu encontrar uma dungeon escondida com a maior dificuldade dentro do Rank D! No fim, enquanto aceitava minha morte, eu ganhei um novo poder!`,
				err:         nil,
				hasContent:  true,
			},
		},
		{
			"error with invalid url",
			args{
				url:        "https://google.com",
				err:        fmt.Errorf("not valid url"),
				hasContent: false,
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			data, err := stub.Info(tt.args.url)
			if !tt.args.hasContent {
				if tt.args.err.Error() != err.Error() {
					t.Errorf("expected err is %v, got %v", tt.args.err, err)
				}
				return
			}
			if data == nil {
				t.Errorf("expected data is not empty, got %v", data)
			}
			if data.Description != tt.args.description {
				t.Errorf("description is not be expected, %v, got %v", tt.args.description, data.Description)
			}

			if data.Title != tt.args.title {
				t.Errorf("title is not be expected, %v, got %v", tt.args.title, data.Title)
			}
			if err != nil {
				t.Errorf("expected err is nil empty, got %v", err)
			}
		})
	}

}
