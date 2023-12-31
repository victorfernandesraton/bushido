package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/victorfernandesraton/bushido/bushido"
	"github.com/victorfernandesraton/bushido/storage"
)

var PageCmd = &cobra.Command{
	Use:              "page [CONTENT_ID] [CHAPTER_ID]",
	Short:            "Get pages from manga",
	Args:             cobra.MinimumNArgs(2),
	TraverseChildren: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		var execSource bushido.Client

		db, err := storage.DatabseFactory()
		if err != nil {
			return err
		}

		contentId, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return errors.New("first arg is not a valid content id")
		}
		chapterId, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return errors.New("second arg is not a valid chapterId id")
		}
		content, err := db.FindById(int(contentId))
		if err != nil {
			return err
		}

		chapter, err := db.FindChapterById(int(chapterId))
		if err != nil {
			return err
		}

		if chapter.Content.ID != content.ID {
			return errors.New("chapter not owner by content")
		} else {
			chapter.Content = content
		}

		sourcesData := Sources()
		execSource, ok := sourcesData[content.Source.ID]
		if !ok {
			return fmt.Errorf(NotFoundSource, content.Source)
		}

		res, err := execSource.Pages(content.ExternalId, chapter.ExternalId)
		if err != nil {
			return err
		}

		persistLocal, err := cmd.Flags().GetBool("dowload")
		if err != nil {
			return err
		}

		if persistLocal {
			err = db.AppendPages(*chapter, res)
			if err != nil {
				return err
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
