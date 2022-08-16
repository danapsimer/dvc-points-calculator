/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile               string
	googleProjectId       string
	googleCredentialsFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dvc-points-calculator",
	Short: "DVC Points Calculator",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dvc-points-calculator.yaml)")
	rootCmd.PersistentFlags().StringVar(&googleProjectId, "google-project-id", "dvc-points-calculator-qa",
		"google project id (default is dvc-points-calculator-qa)")
	rootCmd.PersistentFlags().StringVar(&googleCredentialsFile, "google-credentials-file", "./google-credentials.json",
		"google credentials file (default is ./google-credentials.json)")

	_ = viper.BindPFlag("google.projectId", rootCmd.PersistentFlags().Lookup("google-project-id"))
	_ = viper.BindPFlag("google.credentialsFile", rootCmd.PersistentFlags().Lookup("google-credentials-file"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".dvc-points-calculator" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".dvc-points-calculator")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
