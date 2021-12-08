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
	"strings"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sumo",
	Short: "Sumo Logic CLI",
	Long: ` 
 ___ _   _ _ __ ___   ___  
/ __| | | | '_ ' _ \ / _ \ 
\__ \ |_| | | | | | | (_) |
|___/\__,_|_| |_| |_|\___/ 

sumo is a Sumo Logic CLI and library for building Sumo Logic dashboards and
applications as well as interacting with Sumo Logic platform capabilities.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sumo-cli)")
	rootCmd.PersistentFlags().StringP("access-id", "i", "", "Your Sumo Logic access ID")
	rootCmd.PersistentFlags().StringP("access-key", "k", "", "Your Sumo Logic access key")
	rootCmd.PersistentFlags().StringP("deployment", "r", "us1", "The deployment code for you Sumo Logic instance")

	viper.BindPFlag("access-key", rootCmd.PersistentFlags().Lookup("access-key"))
	viper.BindPFlag("access-id", rootCmd.PersistentFlags().Lookup("access-id"))
	viper.BindPFlag("deployment", rootCmd.PersistentFlags().Lookup("deployment"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".sumo-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".sumo-cli")
		viper.SetConfigType("yaml")
		viper.SetEnvPrefix("SUMO")
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
