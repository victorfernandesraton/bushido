package cmd

import (
	"errors"
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
		var intId, chapterId int64
		var contentId, sourceStr string
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
			if err != nil {
				return err
			}
			sourceStr = info.Source
			contentId = info.ExternalId

		} else {
			sourceStr, err = cmd.Flags().GetString("source")
			if err != nil {
				return err
			}
		}

		if len(args) >= 2 {
			chapterId, err = strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				chapterId = 0
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

		persistLocal, err := cmd.Flags().GetBool("dowload")
		if err != nil {
			return err
		}
		if persistLocal {
			if intId != 0 && chapterId != 0 {
				chapter, err := db.FindChapterById(int(chapterId))
				if err != nil {
					return err
				}

				err = db.AppendPages(int(intId), chapter.ID, sourceStr, res)
				if err != nil {
					return err
				}
			} else {
				return errors.New("not found content or chapter in local storage, See more in chapter add and sync commnad")
			}
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
	PageCmd.Flags().BoolP("dowload", "a", false, "Persist pages in local storage")
	PageCmd.MarkPersistentFlagRequired("source")

}
