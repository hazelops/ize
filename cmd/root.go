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
	"github.com/hazelops/ize/docker"
	"github.com/spf13/cobra"
	"os"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ize",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}



var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
}

var listContainers = &cobra.Command{
	Use:   "list",
	Short: "List Containers",
	Long:  `Lists containers as a test of docker api`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing containers")
		docker.ListContainers()
	},
}

var waypointCmd = &cobra.Command{
	Use:   "waypoint",
	Short: "Run Waypoint",
	Long:  `Runs Waypoint`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Waypoint Command")
		//docker.WaypointInit()
	},
}

var runWaypointInit = &cobra.Command{
	Use:   "init",
	Short: "Run Waypoint init",
	Long:  `Runs Waypoint init`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running Waypoint init")
		docker.WaypointInit()
	},
}


var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Run Terraform Init",
	Long:  `Run Terraform Init`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Terraform Command")
		//docker.TerraformInit()
	},
}

var runTerraformInit = &cobra.Command{
	Use:   "init",
	Short: "Run Terraform Init",
	Long:  `Run Terraform Init`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running Terraform Init")
		docker.TerraformInit()
	},
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ize.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listContainers)
	rootCmd.AddCommand(waypointCmd)
	rootCmd.AddCommand(terraformCmd)

	waypointCmd.AddCommand(runWaypointInit)
	terraformCmd.AddCommand(runTerraformInit)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		//// Find home directory.
		//home, err := os.UserHomeDir()
		//cobra.CheckErr(err)

		// Search config in home directory with name ".ize" (without extension).
		viper.AddConfigPath(".")


		viper.SetConfigName("ize")
		viper.SetConfigType("yaml")
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory")
	}

	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)


	viper.AutomaticEnv() // read in environment variables that match

	//TODO ensure values of the variables are checked for nil before passing down to docker.

	// Global
	viper.SetDefault("ROOT_DIR", cwd)
	viper.SetDefault("INFRA_DIR", fmt.Sprintf("%v/.infra", cwd))
	viper.SetDefault("ENV_DIR", fmt.Sprintf("%v/.infra/env/%v", cwd, viper.Get("ENV")))
	viper.SetDefault("HOME", fmt.Sprintf("%v", home))
	viper.SetDefault("TF_LOG", fmt.Sprintf(""))
	viper.SetDefault("TF_LOG_PATH", fmt.Sprintf("%v/tflog.txt",viper.Get("ENV_DIR")  ))
	viper.SetDefault("TERRAFORM_VERSION", fmt.Sprintf("0.12.29"))


	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

}
