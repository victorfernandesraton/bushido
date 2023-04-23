package cmd

import (
	"github.com/spf13/cobra"
	"github.com/victorfernandesraton/bushido"
	"github.com/victorfernandesraton/bushido/sources/mangalivre"
)

const NotFoundSource = "not found source %v"

var selectedSource string

var RootCmd = &cobra.Command{
	Use:              "bushido [COMAND]",
	Short:            "Bushido is a manga sourece manageer",
	Long:             `Bushido is a manga source manager, a simple way to manager, read , sync and read mangas from diferent sourcers`,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {

	RootCmd.PersistentFlags().StringP("source", "s", selectedSource, "source for content")

}

func Sources() map[string]bushido.Client {
	return map[string]bushido.Client{
		"mangalivre": mangalivre.New(),
	}
}
