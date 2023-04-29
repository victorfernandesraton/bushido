package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var PageCmd = &cobra.Command{
	Use:              "page [CONTENT_ID] [CHAPTER_ID]",
	Short:            "Get pages from manga",
	Args:             cobra.MinimumNArgs(2),
	TraverseChildren: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		sourcesData := Sources()
		selectedSource, err := cmd.Flags().GetString("source")
		if err != nil {
			return err
		}
		execSource, ok := sourcesData[selectedSource]
		if !ok {
			return fmt.Errorf(NotFoundSource, selectedSource)

		}
		res, err := execSource.Pages(args[0], args[1])
		if err != nil {
			return err
		}
		var rows [][]string
		for k, content := range res {
			rows = append(rows, []string{fmt.Sprintf("%d", k), string(content)})
		}
		table := RenderTable([]string{"Number", "Link"}, rows)

		table.Render()

		return nil
	},
}

func init() {

	PageCmd.MarkPersistentFlagRequired("source")

}
