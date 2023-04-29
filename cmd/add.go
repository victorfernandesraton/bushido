package cmd

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/victorfernandesraton/bushido/storage"
)

var AddCmd = &cobra.Command{
	Use:              "add [LINK]",
	Short:            "add manga in local storage",
	Args:             cobra.MinimumNArgs(1),
	TraverseChildren: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		sourcesData := Sources()

		db, err := sql.Open("sqlite3", "sqlite-bushido.db")
		if err != nil {
			panic(err)
		}

		st := storage.New(db)

		if err := st.CreateTables(); err != nil {
			panic(err)
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

		log.Println(res)

		if err := st.Add(*res); err != nil {
			return err
		}

		return nil
	},
}
