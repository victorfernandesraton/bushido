package cmd

import "github.com/spf13/cobra"


var RemoteCmd = &cobra.Command{
	Use:              "remote [COMAND]",
	Short:            "Remote is way to use sources",
	Long:             "Remote is a collection of source related commands such as search in specific manga source or get info using link as origin",
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {

	RemoteCmd.PersistentFlags().StringP("source", "s", selectedSource, "source for content")

}


