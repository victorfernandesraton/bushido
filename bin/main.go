package main

import (
	"fmt"
	"os"

	"github.com/victorfernandesraton/bushido/cmd"
)

func init() {
  cmd.RemoteCmd.AddCommand(cmd.SearchCmd, cmd.InfoCmd, cmd.AddCmd)
  cmd.LocalCmd.AddCommand(cmd.ListCmd)
  cmd.RootCmd.AddCommand(cmd.LocalCmd, cmd.RemoteCmd, cmd.ChapterCmd, cmd.PageCmd, cmd.SyncCmd, cmd.SourceCmd)
}

func main() {

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
