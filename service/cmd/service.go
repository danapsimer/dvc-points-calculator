/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/danapsimer/dvc-points-calculator/api/gin"
	"github.com/danapsimer/dvc-points-calculator/api/goa"
	"github.com/spf13/viper"
	"log"

	"github.com/spf13/cobra"
)

var (
	listenAddresses []string
	engine          string
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "run the service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		switch engine {
		case "goa":
			err = goa.Start()
		case "gin":
			err = gin.Start()
		}
		if err != nil {
			log.Fatalf("error starting engine %s: %s", engine, err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)

	// Here you will define your flags and configuration settings.
	serviceCmd.PersistentFlags().StringArrayVar(&listenAddresses, "listen-addresses", []string{"localhost:8080"},
		"'host:port's to listen on. (default is localhost:8080")
	serviceCmd.PersistentFlags().StringVar(&engine, "engine", "goa",
		"specify the engine to use. can be one of ('gin', 'goa') - (default is goa)")

	_ = viper.BindPFlag("service.listenAddresses", serviceCmd.PersistentFlags().Lookup("listen-addresses"))
	_ = viper.BindPFlag("service.engine", serviceCmd.PersistentFlags().Lookup("engine"))

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
