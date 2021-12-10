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
var diffBuildsCmd = &cobra.Command{
	Use:   "diff-builds [build A] [build B]",
	Short: "Diff objects between two app builds",
	Long: `List all differences between two app builds. This command will compare
all of the folders, dashbaords, panels, saved searches, and variables between two
app builds JSON files and list all of the objects that are created, deleted, and modified,
including what modifications are made.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Fprintf(os.Stderr, "Error: wrong number of arguments. Expects two. Use --help to learn more")
			os.Exit(1)
		}

		appBuild1 := args[0]
		appBuild2 := args[1]

		app := sumoapp.NewApplicationWithPath(appPath)
		baseOverlay := app.NewAppOverlay("base")
		build2DiffOverlay := app.NewAppOverlay("build2Diff")

		baseOverlay.Child = build2DiffOverlay
		build2DiffOverlay.Parent = baseOverlay

		if err := app.ImportToOverlay(appBuild1, baseOverlay); err != nil {
			msg := fmt.Sprintf("Error: could not load %s - %w", appBuild1, err.Error())
			fmt.Fprintf(os.Stderr, msg)
			os.Exit(1)
		}

		if err := app.ImportToOverlay(appBuild2, build2DiffOverlay); err != nil {
			msg := fmt.Sprintf("Error: could not load %s - %w", appBuild1, err.Error())
			fmt.Fprintf(os.Stderr, msg)
			os.Exit(1)
		}

		baseOverlay.Diff(build2DiffOverlay)
	},
}

func init() {
	appCmd.AddCommand(diffBuildsCmd)
}
