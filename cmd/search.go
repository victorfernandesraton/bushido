package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/victorfernandesraton/bushido"
)

var SearchCmd = &cobra.Command{
	Use:              "search [QUERY]",
	Short:            "Search from manga in remote",
	Args:             cobra.MinimumNArgs(1),
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		sourcesData := Sources()
		selectedSource, err := cmd.Flags().GetString("source")
		if err != nil {
			return err
		}

		var contents []bushido.Content
		for _, v := range sourcesData {
			if selectedSource != "" {
				execSource, ok := sourcesData[selectedSource]
				if !ok {
					return fmt.Errorf(NotFoundSource, selectedSource)
				}
				res, err := execSource.Search(args[0])
				if err != nil {
					return err
				}
				contents = append(contents, res...)
			} else {
				res, err := v.Search(args[0])
				if err != nil {
					return err
				}
				contents = append(contents, res...)
			}
		}

		var rows [][]string
		for _, content := range contents {
			rows = append(rows, []string{content.Title, content.Author, content.Link, content.Source})
		}

		table := RenderTable([]string{"Title", "Author", "Link", "source"}, rows)
		table.Render()
		return nil

	},
}
