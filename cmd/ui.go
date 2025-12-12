package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tuannvm/mcpenetes/internal/log"
	"github.com/tuannvm/mcpenetes/internal/ui"
)

// uiCmd represents the ui command
var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Starts the web-based user interface",
	Long:  `Starts a local web server and opens the dashboard in your default browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		noBrowser, _ := cmd.Flags().GetBool("no-browser")

		// Create server
		server := ui.NewServer(port)

		// Open browser in a goroutine (wait a bit for server to start)
		if !noBrowser {
			go func() {
				url := fmt.Sprintf("http://localhost:%d", port)
				_ = openBrowser(url) // reused from search.go (need to verify scope)
			}()
		}

		if err := server.Start(); err != nil {
			log.Fatal("Server failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().IntP("port", "p", 3000, "Port to run the UI server on")
	uiCmd.Flags().Bool("no-browser", false, "Do not open the default browser automatically")
}
