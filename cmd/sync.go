package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/victorfernandesraton/bushido"
)

var SyncCmd = &cobra.Command{
	Use:              "sync [ID]",
	Short:            "Sync manga chapters from remote to local storage",
	Args:             cobra.MinimumNArgs(0),
	TraverseChildren: true,

	RunE: func(cmd *cobra.Command, args []string) error {

		var intId int64
		var contents []bushido.Content

		recursiveSearch, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}

		syncAll, err := cmd.Flags().GetBool("sync-all")
		if err != nil {
			return err
		}

		db, err := DatabseFactory()
		if err != nil {
			return err
		}

		if len(args) >= 1 || !syncAll {
			intId, err = strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			info, err := db.FindById(int(intId))
			if err != nil {
				return err
			}

			contents = append(contents, *info)
		} else {
			contents, err = db.ListByName("")
			if err != nil {
				return err
			}
		}

		log.Println("contents to sync ", len(contents))

		sourcesData := Sources()

		for _, c := range contents {
			execSource, ok := sourcesData[c.Source]
			if !ok {
				return fmt.Errorf(NotFoundSource, selectedSource)
			}

			chapters, err := execSource.Chapters(c.Link, recursiveSearch)
			log.Println(fmt.Printf("%v chapters to sync for %v from %v", len(chapters), c.Title, c.Source))
			if err != nil {
				return err
			}

			if err := db.AppendChapter(int(intId), chapters); err != nil {
				return err
			}

		}

		return nil
	},
}

func init() {
	SyncCmd.Flags().BoolP("recursive", "r", false, "Search chapters with recursion and list all")
	SyncCmd.Flags().Bool("sync-all", false, "sync all local data")
}
