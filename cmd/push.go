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
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sumologic.com/sumo-cli/sumoapp"
)

var (
	appDestinationParent    string
	appDestinationName      string
	appDestinationOverwrite bool
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push an application build to your Sumo Logic account",
	Args:  cobra.MinimumNArgs(1),
	Long: `Push an application build (a json file, see 'sumo app build --help' for more information) to
your Sumo Logic organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		var apiURL string
		var buildPath string

		rootFolder := sumoapp.NewFolder()

		accessId := viper.GetString("access-id")
		accessKey := viper.GetString("access-key")
		region := viper.GetString("deployment")

		if region == "us1" {
			apiURL = "https://api.sumologic.com/api"
		} else {
			apiURL = fmt.Sprintf("https://api.%s.sumologic.com/api", region)
		}

		switch len(args) {
		case 0:
			buildPath = "./build.json"
		case 1:
			buildPath = args[0]
		default:
			fmt.Fprintf(os.Stderr, "Error: too many arguments. Expects none or one. Use --help to learn more")
			os.Exit(1)
		}

		//Read in build file
		data, err := os.ReadFile(buildPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		if err := json.Unmarshal(data, &rootFolder); err != nil {
			if jsonErr, ok := err.(*json.SyntaxError); ok {
				problemPart := data[jsonErr.Offset-10 : jsonErr.Offset+10]
				err = fmt.Errorf("%w ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
				fmt.Fprintf(os.Stderr, err.Error())
			} else {
				fmt.Fprintf(os.Stderr, err.Error())
			}

			os.Exit(1)
		}

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

		should_overwrite, _ := cmd.Flags().GetBool("overwrite")

		if err := rootFolder.Upload(&client, appDestinationParent, should_overwrite); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	appCmd.AddCommand(pushCmd)

	pushCmd.PersistentFlags().StringVarP(&appDestinationParent, "parent-folder", "d", "", "ID of the folder to put the application into")
	pushCmd.PersistentFlags().BoolP("overwrite", "w", false, "Whether to overwrite an existing destination folder")
}
