package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/victorfernandesraton/bushido/web"
)

var ServerCmd = &cobra.Command{
	Use:              "server",
	Short:            "Start server for self-hosted platform",
	Args:             cobra.MinimumNArgs(0),
	TraverseChildren: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		httpServer := web.NewWebServer(3000)
		log.Printf("Server start runing in %s", httpServer.Addr)
		return httpServer.ListenAndServe()

	},
}
