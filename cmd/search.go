package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/tuannvm/mcpenetes/internal/config"
	"github.com/tuannvm/mcpenetes/internal/log"
	"github.com/tuannvm/mcpenetes/internal/registry"
	"github.com/tuannvm/mcpenetes/internal/search"
)

// ServerInfo represents information about an MCP server
type ServerInfo struct {
	Name          string
	Description   string
	RepositoryURL string
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [server-id]",
	Short: "Interactive fuzzy search for MCP versions and apply them",
	Long: `Provides an interactive fuzzy search interface to find and select MCP versions from configured registries.
If a server ID is provided as an argument, it will directly use that server without prompting.
After selection, the server is added to the local mcp.json configuration file.
Note: Since the registry may not provide specific execution commands, a default 'npx' configuration is added.
You should review 'mcp.json' to ensure the command and arguments are correct.

By default, search results are cached to improve performance. Use the --refresh flag to force a refresh
of the cache and fetch the latest data from the registries.`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Allow 0 or 1 argument
		if len(args) > 1 {
			return errors.New("accepts at most one argument: the server ID to use")
		}
		if len(args) == 1 && args[0] == "" {
			return errors.New("server ID cannot be empty if provided")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Get the refresh flag value
		forceRefresh, _ := cmd.Flags().GetBool("refresh")
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatal("Error loading config: %v", err)
		}

		var serverID string
		var selectedServer *registry.ServerData

		// Direct server selection if argument provided
		if len(args) == 1 {
			serverID = args[0]
			log.Info("Using provided server ID: %s", serverID)
			// We don't have server data if passed directly, unless we fetch?
			// For simplicity, we just use the ID.
		} else {
			// Interactive selection mode
			log.Info("Starting interactive search...")

			if len(cfg.Registries) == 0 {
				log.Warn("No registries configured. Use 'mcpetes add registry <n> <url>' to add one.")
				return
			}

			// Start spinner
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Suffix = " Fetching available MCPs..."
			s.Start()

			var serverInfos []ServerInfo
			var displayOptions []string
			serverMap := make(map[string]registry.ServerData)

			for _, reg := range cfg.Registries {
				servers, err := registry.FetchMCPServersWithCache(reg.URL, forceRefresh)
				if err != nil {
					log.Warn("Error fetching from registry %s: %v", reg.URL, err)
					continue
				}

				for _, server := range servers {
					info := ServerInfo{
						Name:          server.Name,
						Description:   server.Description,
						RepositoryURL: server.RepositoryURL,
					}
					serverInfos = append(serverInfos, info)

					// Store original server data for later use
					serverMap[server.Name] = server

					// Create display option string
					displayText := server.Name
					if server.Description != "" {
						displayText = fmt.Sprintf("%s: %s", server.Name, server.Description)
					}
					displayOptions = append(displayOptions, displayText)
				}
			}

			s.Stop()

			if len(serverInfos) == 0 {
				log.Warn("No MCP servers found in any registry")
				return
			}

			var selectedOption string
			prompt := &survey.Select{
				Message: "Select MCP server:",
				Options: displayOptions,
			}

			err = survey.AskOne(prompt, &selectedOption)
			if err != nil {
				log.Fatal("Error during selection: %v", err)
				return
			}

			// Find the index/ID
			for i, opt := range displayOptions {
				if opt == selectedOption {
					serverID = serverInfos[i].Name
					// Get full data
					if data, ok := serverMap[serverID]; ok {
						selectedServer = &data
					}
					break
				}
			}
		}

		if serverID == "" {
			log.Fatal("No server selected")
		}

		log.Info("Selected MCP: %s", serverID)

		// Ask to open repo if available
		if selectedServer != nil && selectedServer.RepositoryURL != "" {
			var openRepo bool
			confirmPrompt := &survey.Confirm{
				Message: fmt.Sprintf("Would you like to open the repository URL (%s) in your browser?", selectedServer.RepositoryURL),
				Default: true,
			}
			_ = survey.AskOne(confirmPrompt, &openRepo) // Ignore error, optional
			if openRepo {
				_ = openBrowser(selectedServer.RepositoryURL)
			}
		}

		// Add to config.yaml (legacy/tracking)
		cfg.MCPs = append(cfg.MCPs, serverID)
		if err := config.SaveConfig(cfg); err != nil {
			log.Fatal("Error saving config: %v", err)
		}

		// Add to mcp.json (actual configuration)
		log.Info("Adding configuration for %s to mcp.json...", serverID)
		err = search.AddServerToMCPConfig(serverID, selectedServer)
		if err != nil {
			log.Error("Failed to add to mcp.json: %v", err)
		} else {
			log.Success("Successfully added %s to mcp.json.", serverID)
			log.Info("Note: A default 'npx' command was configured. Please check 'mcp.json' if you need to adjust arguments or environment variables.")
			log.Info("Run 'mcpenetes apply' to install this server to your clients.")
		}
	},
}

// openBrowser opens the specified URL in the default browser
func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}

func init() {
	rootCmd.AddCommand(searchCmd)
	
	// Add a flag to force cache refresh
	searchCmd.Flags().BoolP("refresh", "r", false, "Force a refresh of the cache and fetch the latest data from registries")
}
