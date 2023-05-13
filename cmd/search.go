package cmd

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"github.com/victorfernandesraton/bushido"
)

func asyncSearch(wg *sync.WaitGroup, client bushido.Client , query string ,result chan<- []bushido.Content) {
  defer wg.Done()
  res, err := client.Search(query)
  if err != nil {
    panic(err)
  }
  result <- res
}

var SearchCmd = &cobra.Command{
	Use:              "search [QUERY]",
	Short:            "Search from manga in remote",
	Args:             cobra.MinimumNArgs(1),
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		sourcesData := Sources()
		selectedSource, err := cmd.Flags().GetString("source")
		if err != nil {
			return err
		}

		var contents []bushido.Content
 		if selectedSource != "" {
		  execSource, ok := sourcesData[selectedSource]
			if !ok {
				return fmt.Errorf(NotFoundSource, selectedSource)
			}
			res, err := execSource.Search(args[0])
			if err != nil {
				return err
			}
			contents = append(contents, res...)
		} else {
      resultch := make(chan []bushido.Content)
      var wg sync.WaitGroup
      wg.Add(len(sourcesData))

      for _, v := range sourcesData {
        go asyncSearch(&wg, v, args[0], resultch )
			}

      go func() {
		    wg.Wait()
		    close(resultch)
	    }()

      for res :=range resultch {
        contents = append(contents, res...)
      }
		}

		var rows [][]string
		for _, content := range contents {
			rows = append(rows, []string{content.Title, content.Author, content.Link, content.Source})
		}

		table := RenderTable([]string{"Title", "Author", "Link", "source"}, rows)
		table.Render()
		return nil

	},
}
