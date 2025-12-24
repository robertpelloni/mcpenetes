package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tuannvm/mcpenetes/internal/doctor"
	"github.com/tuannvm/mcpenetes/internal/log"
)

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system health and prerequisites",
	Long:  `Runs a series of checks to ensure mcpenetes and its dependencies are configured correctly.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Running system health checks...\n")

		results := doctor.RunChecks()
		hasError := false

		for _, res := range results {
			switch res.Status {
			case "ok":
				fmt.Printf("✅ %s: %s\n", res.Name, res.Message)
			case "warning":
				fmt.Printf("⚠️  %s: %s\n", res.Name, res.Message)
			case "error":
				fmt.Printf("❌ %s: %s\n", res.Name, res.Message)
				hasError = true
			}
		}

		fmt.Println()
		if hasError {
			log.Warn("Some checks failed. Please address the issues above.")
			os.Exit(1)
		} else {
			log.Success("All checks passed!")
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
