package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var SyncCmd = &cobra.Command{
	Use:              "sync [ID]",
	Short:            "sync manga chapter from id",
	Args:             cobra.MinimumNArgs(0),
	TraverseChildren: true,

	RunE: func(cmd *cobra.Command, args []string) error {

		var intId int64

		recursiveSearch, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}

		db, err := DatabseFactory()
		if err != nil {
			return err
		}

		if len(args) >= 1 {
			intId, err = strconv.ParseInt(args[0], 10, 64)
			fmt.Println(intId)
			if err != nil {
				return err
			}
		}

		info, err := db.FindById(int(intId))
		if err != nil {
			return err
		}

		if info == nil {
			return errors.New("content not exist in base")
		}

		sourcesData := Sources()
		execSource, ok := sourcesData[info.Source]
		fmt.Println(info.BasicContent.Source)
		if !ok {
			return fmt.Errorf(NotFoundSource, selectedSource)

		}

		res, err := execSource.Chapters(info.Link, recursiveSearch)
		if err != nil {
			return err
		}
		if err := db.AppendChapter(int(intId), res); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	SyncCmd.Flags().BoolP("recursive", "r", false, "Search chapters with recursion and list all")
}
