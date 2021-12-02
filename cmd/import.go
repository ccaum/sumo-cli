/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	appStream string
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a sumo application",
	Long: `Import an existing folder or other set of resources. The resources
will be broken into components (i.e. dashboards, folders, panels, variables).
By default, the resources will be put into the 'upstream' application stream.
You can override this behavior using the --app-stream parameter.`,

	Run: func(cmd *cobra.Command, args []string) {
		var filePath string

		switch len(args) {
		case 0:
			fmt.Fprintf(os.Stderr, "Error: not enough arguments. Provide the json file to import. See https://help.sumologic.com/01Start-Here/Library/Export-and-Import-Content-in-the-Library#export-content-in-the-library for more information")
			os.Exit(1)
		case 1:
			filePath = args[0]
		default:
			fmt.Fprintf(os.Stderr, "Error: too many arguments. Expects a single argument with the path to the json file to import. See https://help.sumologic.com/01Start-Here/Library/Export-and-Import-Content-in-the-Library#export-content-in-the-library for more information")
		}

		if err := sumoapp.Import(filePath, appStream); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.PersistentFlags().StringVarP(&appStream, "app-stream", "s", "upstream", "Which app stream to import to")
}
