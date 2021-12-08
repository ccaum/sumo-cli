package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
	Short: "Test loading an appstream",
	Long:  `Test loading an appstream`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			fmt.Print("Usage: debug load app-stream type key")
			os.Exit(1)
		}

		appStream := args[0]
		thetype := args[1]
		key := args[2]

		app := sumoapp.NewApplicationWithPath(appPath)
		if err := app.LoadAppStreams(); err != nil {
			msg := fmt.Errorf("Unable to load app streams: %w", err)
			fmt.Println(msg)
			os.Exit(1)
		}

		stream, err := app.FindAppStream(appStream)
		if err != nil {
			msg := fmt.Errorf("Unable to load app stream %s: %w", appStream, err)
			fmt.Println(msg)
			os.Exit(1)
		}

		switch thetype {
		case "variable":
			fmt.Print(stream.Variables[key])
		case "panel":
			fmt.Print(stream.Panels[key])
		case "saved-search":
			fmt.Print(stream.SavedSearches[key])
		case "dashboards":
			fmt.Print(stream.Dashboards[key])
		case "folder":
			fmt.Print(stream.Folders[key])
		}
	},
}

func init() {
	debugCmd.AddCommand(debugLoadCmd)
	rootCmd.AddCommand(debugCmd)

	debugCmd.PersistentFlags().StringVarP(&appPath, "app-path", "p", ".", "The path to the application")
}
