/*
Copyright © 2021 Carl Caum <carl@carlcaum.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sumologic.com/sumo-cli/sumoapp"
)

// diffOverlaysCmd represents the diff-overlays command
var diffOverlaysCmd = &cobra.Command{
	Use:   "diff-overlays [overlay name] [overlay name]",
	Short: "Diff objects between two app overlays",
	Long: `List all differences between two app overlays. This command will compare
all of the folders, dashbaords, panels, saved searches, and variables between two
app overlays and list all of the objects that are created, deleted, and modified,
including what modifications are made.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Fprintf(os.Stderr, "Error: wrong number of arguments. Expects two. Use --help to learn more")
			os.Exit(1)
		}

		appOverlay1 := args[0]
		appOverlay2 := args[1]

		app := sumoapp.NewApplicationWithPath(appPath)
		if err := app.LoadAppOverlays(); err != nil {
			msg := fmt.Errorf("Unable to load app overlays: %w", err)
			fmt.Println(msg)
			os.Exit(1)
		}

		overlay1, err := app.FindAppOverlay(appOverlay1)
		if err != nil {
			msg := fmt.Errorf("Unable to load app overlay %s: %w", appOverlay1, err)
			fmt.Println(msg)
			os.Exit(1)
		}

		overlay2, err := app.FindAppOverlay(appOverlay2)
		if err != nil {
			msg := fmt.Errorf("Unable to load app overlay %s: %w", appOverlay2, err)
			fmt.Println(msg)
			os.Exit(1)
		}

		overlay1.Diff(overlay2)
	},
}

func init() {
	appCmd.AddCommand(diffOverlaysCmd)
}
