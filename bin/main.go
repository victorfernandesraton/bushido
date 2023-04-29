package main

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/victorfernandesraton/bushido/cmd"
)

func init() {
	cmd.RootCmd.AddCommand(cmd.InfoCmd, cmd.ChapterCmd, cmd.SearchCmd, cmd.PageCmd, cmd.AddCmd, cmd.ListCmd, cmd.SyncCmd)
}

func main() {

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
