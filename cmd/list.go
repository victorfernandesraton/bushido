package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/victorfernandesraton/bushido/bushido"
	"github.com/victorfernandesraton/bushido/storage"
)

var ListCmd = &cobra.Command{
	Use:              "list [LINK]",
	Short:            "List mangas storage in local",
	Args:             cobra.MinimumNArgs(0),
	TraverseChildren: true,

	RunE: func(cmd *cobra.Command, args []string) error {

		var searchText string
		var res []bushido.Content
		var intId int64

		db, err := storage.DatabseFactory()
		if err != nil {
			return err
		}

		if len(args) >= 1 {
			intId, err = strconv.ParseInt(args[0], 10, 64)
			fmt.Println(intId)
			if err != nil {
				searchText = args[0]
			}
		}

		if intId != 0 {
			data, err := db.FindById(int(intId))
			res = append(res, *data)
			if err != nil {
				return err
			}

		} else {
			res, err = db.ListByName(searchText)
			if err != nil {
				return err
			}
		}

		var rows [][]string
		for _, content := range res {
			rows = append(rows, []string{strconv.Itoa(content.ID), content.Title, content.Author, content.Link, content.Source.ID})
		}

		table := RenderTable([]string{"ID", "Title", "Author", "Link", "source"}, rows)

		table.Render()

		return nil
	},
}
