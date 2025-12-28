package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/tuannvm/mcpenetes/internal/proxy"
)

var proxyServerID string

var proxyCmd = &cobra.Command{
	Use:   "proxy [flags] -- [command] [args...]",
	Short: "Run a command wrapped in the mcpenetes logging proxy",
	Long: `Executes the specified command while capturing stderr to a log file.
This allows mcpenetes to provide logs for MCP servers running inside other clients.

Example:
  mcpenetes proxy --server-id my-server -- npx -y @modelcontextprotocol/server-filesystem /path/to/files

  Note: Use '--' to separate the proxy flags from the command to be executed.`,
	DisableFlagParsing: false, // We want to parse our flags
	RunE: func(c *cobra.Command, args []string) error {
		// Args contains everything after the flags if we use standard parsing
		// Example: mcpenetes proxy --server-id foo npx -y bar
		// Cobra parses --server-id foo.
		// Args is ["npx", "-y", "bar"].

		if len(args) < 1 {
			return fmt.Errorf("requires at least a command to run")
		}

		opts := proxy.Options{
			Command:  args[0],
			Args:     args[1:],
			ServerID: proxyServerID,
		}

		if err := proxy.Run(opts); err != nil {
			// If the command fails, we exit with its code if possible, or 1
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			return err
		}
		return nil
	},
}

func init() {
	proxyCmd.Flags().StringVar(&proxyServerID, "server-id", "", "Unique ID for the server (for log naming)")

	// We enable FParseAllWhitelisted so that flags after the first non-flag arg are NOT treated as flags for us.
	// But actually, `proxy [flags] [cmd] [args...]` works fine with standard Interspersed=true (default),
	// unless [cmd] starts with `-`.
	// To be safe, users should use `--` separator if their command starts with `-`.
	// Or we can disable interspersed parsing.
	proxyCmd.Flags().SetInterspersed(false)

	rootCmd.AddCommand(proxyCmd)
}
