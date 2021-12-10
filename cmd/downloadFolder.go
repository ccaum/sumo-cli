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
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sumologic.com/sumo-cli/sumoapp"
)

var (
	downloadDestination string
)

// downloadCmd represents the push command
var downloadFoldersCmd = &cobra.Command{
	Use:   "download-folder",
	Short: "Download an application folder from your Sumo Logic account",
	Args:  cobra.MinimumNArgs(1),
	Long:  `Download an application folder from your Sumo Logic account.`,
	Run: func(cmd *cobra.Command, args []string) {
		var apiURL string

		rootFolder := sumoapp.NewFolder()

		accessId := viper.GetString("access-id")
		accessKey := viper.GetString("access-key")
		region := viper.GetString("deployment")

		if region == "us1" {
			apiURL = "https://api.sumologic.com/api"
		} else {
			apiURL = fmt.Sprintf("https://api.%s.sumologic.com/api", region)
		}

		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "Error: Expects one argument. Use --help to learn more")
			os.Exit(1)
		}

		rootFolder.Id = args[0]

		client := sumoapp.APIClient{
			Cfg: &sumoapp.Configuration{
				Authentication: sumoapp.BasicAuth{
					AccessId:  accessId,
					AccessKey: accessKey,
				},
				BasePath:   apiURL,
				HTTPClient: &http.Client{},
			},
		}

		fileBytes, err := rootFolder.Download(&client)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Println("D: ", downloadDestination)
		if downloadDestination == "-" {
			fmt.Println(string(fileBytes))
		} else {
			os.WriteFile(downloadDestination, fileBytes, 0644)
		}
	},
}

func init() {
	appCmd.AddCommand(downloadFoldersCmd)

	downloadFoldersCmd.PersistentFlags().StringVarP(&downloadDestination, "output-file", "o", "-", "File to save the folder's output file to")
}
