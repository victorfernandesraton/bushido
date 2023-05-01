package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var SourceCmd = &cobra.Command{
  Use:              "source",
	Short:            "list active sources",
	Args:             cobra.MinimumNArgs(0),
	TraverseChildren: true,

  RunE: func(cmd *cobra.Command, args []string) error {
    sourcesData := Sources()
    sources := make([]string, 0 , len(sourcesData))
    for k := range sourcesData {
      sources = append(sources, k)
    }

   	var rows [][]string
		for k, content := range sources {
			rows = append(rows, []string{fmt.Sprintf("%d", k), content})
		}
		table := RenderTable([]string{"Number", "Source"}, rows)

		table.Render()

		return nil 
  },
}	
