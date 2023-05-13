package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/victorfernandesraton/bushido/storage"
)
var AddCmd = &cobra.Command{

	Use:              "add [LINK]",
	Short:            "Add manga in local storage",
	Args:             cobra.MinimumNArgs(1),
	TraverseChildren: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		sourcesData := Sources()
		db, err := storage.DatabseFactory()

		if err != nil {
			return err
		}

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

		if err := db.Add(*res); err != nil {
			return err
		}

    log.Println(fmt.Sprintf("Manga %v has been added in local storage", res.Title))

		return nil
	},
}
