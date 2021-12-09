package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"sumologic.com/sumo-cli/sumoapp"
)

// buildCmd represents the build command
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Test command",
	Long:  `Use a subcommand to test a feature`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("debug called")
	},
}

var debugLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Test loading an appoverlay",
	Long:  `Test loading an appoverlay`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			fmt.Print("Usage: debug load app-overlay type key")
			os.Exit(1)
		}

		appOverlay := args[0]
		thetype := args[1]
		key := args[2]

		app := sumoapp.NewApplicationWithPath(appPath)
		if err := app.LoadAppOverlays(); err != nil {
			msg := fmt.Errorf("Unable to load app overlays: %w", err)
			fmt.Println(msg)
			os.Exit(1)
		}

		overlay, err := app.FindAppOverlay(appOverlay)
		if err != nil {
			msg := fmt.Errorf("Unable to load app overlay %s: %w", appOverlay, err)
			fmt.Println(msg)
			os.Exit(1)
		}

		switch thetype {
		case "variable":
			WriteYamlObject(overlay.Variables[key])
		case "panel":
			WriteYamlObject(overlay.Panels[key])
		case "saved-search":
			WriteYamlObject(overlay.SavedSearches[key])
		case "dashboards":
			WriteYamlObject(overlay.Dashboards[key])
		case "folder":
			WriteYamlObject(overlay.Folders[key])
		}
	},
}

func WriteYamlObject(object interface{}) {
	p, err := yaml.Marshal(&object)
	if err != nil {
		msg := fmt.Errorf("Unable to marshall: %w", err)
		fmt.Print(msg)
	}
	fmt.Print(string(p))
}

func init() {
	debugCmd.AddCommand(debugLoadCmd)
	rootCmd.AddCommand(debugCmd)

	debugCmd.PersistentFlags().StringVarP(&appPath, "app-path", "p", ".", "The path to the application")
}
