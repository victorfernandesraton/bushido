package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/victorfernandesraton/bushido"
)

var ChapterCmd = &cobra.Command{
	Use:              "chapter [LINK]",
	Short:            "Get chapters from manga link",
	Args:             cobra.MinimumNArgs(1),
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var res []bushido.Chapter
		var intId int64
		var link string
		var sourceStr string

		sourcesData := Sources()
		db, err := DatabseFactory()
		if err != nil {
			return err
		}

		if len(args) >= 1 {
			intId, err = strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				link = args[0]
				sourceStr, err = cmd.Flags().GetString("source")
				if err != nil {
					return err
				}
			} else {
				info, err := db.FindById(int(intId))
				if err != nil {
					return err
				}
				link = info.Link
				sourceStr = info.Source
			}
		}

		recursiveSearch, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}
		localOnly, err := cmd.Flags().GetBool("local-only")
		if err != nil {
			return err
		}

		execSource, ok := sourcesData[sourceStr]
		if !ok {
			return fmt.Errorf(NotFoundSource, selectedSource)

		}

		if localOnly && intId != 0 {
			res, err = db.ListChaptersByContentId(int(intId))
			if err != nil {
				return err
			}
		} else {
			res, err = execSource.Chapters(link, recursiveSearch)
			if err != nil {
				return err
			}
		}

		var rows [][]string
		for _, content := range res {
			rows = append(rows, []string{strconv.Itoa(content.ID), content.Title, content.ExternalId, content.Link, content.Content.ExternalId})
		}

		table := RenderTable([]string{"ID", "Title", "id", "link", "contend_id"}, rows)
		table.Render()
		return nil
	},
}

func init() {
	ChapterCmd.Flags().BoolP("recursive", "r", false, "Search chapters with recursion and list all")
	ChapterCmd.Flags().Bool("local-only", false, "Search chapters lcoal only")

}
