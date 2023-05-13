package cmd

import "github.com/spf13/cobra"

var LocalCmd = &cobra.Command{
	Use:              "local [COMAND]",
	Short:            "Local is a collection of commands with depend of local storage",
	Long:             "Local is a collection of commands with we use to manager and consume local data",
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {

	},
}





