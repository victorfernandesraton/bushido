package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/victorfernandesraton/bushido"
	"github.com/victorfernandesraton/bushido/sources/mangalivre"
)

const NotFoundSource = "not found source %v"

var selectedSource string
var sourcesData map[string]bushido.Client

var rootCmd = &cobra.Command{
	Use:              "bushido [COMAND]",
	Short:            "Bushido is a manga sourece manageer",
	Long:             `Bushido is a manga source manager, a simple way to manager, read , sync and read mangas from diferent sourcers`,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var searchCmd = &cobra.Command{
	Use:              "search [QUERY]",
	Short:            "Search from manga in remote",
	Args:             cobra.MinimumNArgs(1),
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		selectedSource, err := cmd.Flags().GetString("source")
		if err != nil {
			return err
		}

		var contents []bushido.Content
		for _, v := range sourcesData {
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
				res, err := v.Search(args[0])
				if err != nil {
					return err
				}
				contents = append(contents, res...)
			}
		}
		fmt.Println(contents)
		return nil

	},
}

var infoCmd = &cobra.Command{
	Use:              "info [LINK]",
	Short:            "Get mangainfo from source",
	Args:             cobra.MinimumNArgs(1),
	TraverseChildren: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		selectedSource, err := cmd.Flags().GetString("source")
		if err != nil {
			return err
		}
		execSource, ok := sourcesData[selectedSource]
		if !ok {
			return fmt.Errorf(NotFoundSource, selectedSource)

		}
		res, err := execSource.Info(args[0])
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}

func init() {
	sourcesData = map[string]bushido.Client{
		"mangalivre": mangalivre.New(),
	}

	rootCmd.PersistentFlags().StringP("source", "s", selectedSource, "source for content")
	infoCmd.MarkPersistentFlagRequired("source")
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(infoCmd)

}
func main() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
