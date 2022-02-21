/*
Copyright © 2022 Rodion Lim <rodion.lim@hotmail.com>

*/
package cmd

import (
	"io"
	"os"

	crs "github.com/rodionlim/carousell/library/carousell"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "carousell",
	Short: "Carousell is a CLI tool that fetches Carousell listings and notifies user about new listings",
	Long: `
Carousell is a CLI library for Go that fetches Carousell listings.
This application integrates with Slack for notifications.
To enable slack notifications, the application expects that
an environment variable SLACK_ACCESS_TOKEN is set.`,
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
	rootCmd.PersistentFlags().BoolP("recent", "r", false, "Search recent listings")
	rootCmd.PersistentFlags().IntP("price-floor", "f", 0, "Minimum price of listing")
	rootCmd.PersistentFlags().IntP("price-ceil", "c", 0, "Maximum price of listing")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose mode with logging")
}

func setVerbosity(cmd *cobra.Command) {
	verbose, _ := rootCmd.PersistentFlags().GetBool("verbose")
	if !verbose {
		logrus.SetOutput(io.Discard)
	}
}

func getOpts(cmd *cobra.Command, args []string) []func(r *crs.Req) {
	var opts []func(r *crs.Req)
	for _, arg := range args {
		opts = append(opts, crs.WithSearch(arg))
	}

	recentFlag, err := cmd.Root().PersistentFlags().GetBool("recent")
	if recentFlag && err == nil {
		opts = append(opts, crs.WithRecent)
	}

	floor, err := cmd.Root().PersistentFlags().GetInt("price-floor")
	if err == nil && floor != 0 {
		opts = append(opts, crs.WithPriceFloor(floor))
	}

	ceil, err := cmd.Root().PersistentFlags().GetInt("price-ceil")
	if err == nil && ceil != 0 {
		opts = append(opts, crs.WithPriceCeil(ceil))
	}
	return opts
}
