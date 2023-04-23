package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var InfoCmd = &cobra.Command{
	Use:              "info [LINK]",
	Short:            "Get mangainfo from source",
	Args:             cobra.MinimumNArgs(1),
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
		res, err := execSource.Info(args[0])
		if err != nil {
			return err
		}
		table := RenderTable([]string{"Title", "Author", "Link", "source"}, [][]string{
			{res.Title, res.Author, res.Link, res.Source},
		})

		table.Render()

		return nil
	},
}

func init() {

	InfoCmd.MarkPersistentFlagRequired("source")

}
