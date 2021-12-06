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

var (
	outputFile string
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Compile a single application JSON artifact",
	Long: `Compiles all the application streams into a single JSON
file that can be imported into Sumo Logic's Continuous Intelligence Platform`,
	Run: func(cmd *cobra.Command, args []string) {
		var path string

		switch len(args) {
		case 0:
			path = appPath
		case 1:
			path = args[0]
		default:
			fmt.Fprintf(os.Stderr, "Error: too many arguments. Expects none or one. Use --help to learn more")
			os.Exit(1)
		}

		app := sumoapp.NewApplicationWithPath(path)
		if err := app.LoadAppStreams(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err)
			os.Exit(1)
		}

		jsonString, err := app.ToJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err)
			os.Exit(1)
		}

		if outputFile == "" {
			fmt.Println(string(jsonString))
		} else {
			err := os.WriteFile(outputFile, jsonString, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Unable to write to file %s. %w", outputFile, err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	appCmd.AddCommand(buildCmd)

	buildCmd.PersistentFlags().StringVarP(&outputFile, "output-file", "o", "", "Output file containing the compiled JSON")
}
