package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/victorfernandesraton/bushido/storage"
)

var InfoCmd = &cobra.Command{
	Use:              "info [LINK]",
	Short:            "Get mangainfo fromb local or source",
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

		db, err := storage.DatabseFactory()
		if err != nil {
			return err
		}

		res, err := db.FindByLink(args[0])
		if err != nil {
			if err.Error() != "content not found" {
				return err
			}
			res, err = execSource.Info(args[0])
			if err != nil {
				return err
			}
		}

		table := RenderTable([]string{"ID", "Title", "Author", "Link", "source"}, [][]string{
			{strconv.Itoa(res.ID), res.Title, res.Author, res.Link, res.Source},
		})

		table.Render()

		return nil
	},
}

func init() {
  InfoCmd.MarkPersistentFlagRequired("source")
}
