package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:              "add [LINK]",
	Short:            "add manga in local storage",
	Args:             cobra.MinimumNArgs(1),
	TraverseChildren: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		sourcesData := Sources()

		db, err := DatabseFactory()
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

		return nil
	},
}
