package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/victorfernandesraton/bushido"
)

var PageCmd = &cobra.Command{
	Use:              "page [CONTENT_ID] [CHAPTER_ID]",
	Short:            "Get pages from manga",
	Args:             cobra.MinimumNArgs(2),
	TraverseChildren: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		var intId int64
		var contentId string
		var sourceStr string
		var execSource bushido.Client

		db, err := DatabseFactory()
		if err != nil {
			return err
		}

		if len(args) >= 1 {
			intId, err = strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				contentId = args[0]
			}
		}

		if intId != 0 {
			info, err := db.FindById(int(intId))
			sourceStr = info.Source
			contentId = info.ExternalId
			if err != nil {
				return err
			}

		} else {
			sourceStr, err = cmd.Flags().GetString("source")
			if err != nil {
				return err
			}
		}
		sourcesData := Sources()
		execSource, ok := sourcesData[sourceStr]
		if !ok {
			return fmt.Errorf(NotFoundSource, sourceStr)

		}
		res, err := execSource.Pages(contentId, args[1])
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
