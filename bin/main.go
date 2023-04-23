package main

import (
	"fmt"
	"os"

	"github.com/victorfernandesraton/bushido/cmd"
)

func init() {
	cmd.RootCmd.AddCommand(cmd.InfoCmd, cmd.ChapterCmd, cmd.SearchCmd)
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
