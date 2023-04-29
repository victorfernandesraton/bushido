package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ChapterCmd = &cobra.Command{
	Use:              "chapter [LINK]",
	Short:            "Get chapters from manga link",
	Args:             cobra.MinimumNArgs(1),
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		sourcesData := Sources()
		selectedSource, err := cmd.Flags().GetString("source")
		if err != nil {
			return err
		}
		recursiveSearch, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}

		execSource, ok := sourcesData[selectedSource]
		if !ok {
			return fmt.Errorf(NotFoundSource, selectedSource)

		}
		res, err := execSource.Chapters(args[0], recursiveSearch)
		if err != nil {
			return err
		}

		var rows [][]string
		for _, content := range res {
			rows = append(rows, []string{content.Title, content.ExternalId, content.Link, content.Content.ExternalId})
		}

		table := RenderTable([]string{"Title", "id", "link", "contend_id"}, rows)
		table.Render()
		return nil
	},
}

func init() {
	ChapterCmd.Flags().BoolP("recursive", "r", false, "Search chapters with recursion and list all")
}
