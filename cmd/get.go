/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	crs "github.com/rodionlim/carousell/library/carousell"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [search terms...]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Fetches carousell listings",
	Long: `
Fetches carousell listings. Arguments supplied to this command
are search terms in Carousell, and at least one search term should be provided.

Flags can be used to modify the search behaviour, such as specifying the 
-r recent flag, to include only recent listings etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := getOpts(cmd, args)
		req := crs.NewReq(opts...)
		listings, err := req.Get()
		if err != nil {
			fmt.Println("Something unexpected happened")
			os.Exit(1)
		}

		// If user specifies that they want a summarized version of the output
		shortFlag, err := cmd.Flags().GetBool("shorthand")
		if err == nil && shortFlag {
			for _, listing := range listings {
				listing.Print()
				fmt.Println()
			}
			return
		}

		fmt.Println("Obtained", len(listings), "listings")
		fmt.Printf("%+v\n", listings)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getCmd.Flags().BoolP("shorthand", "s", false, "Display listings output in summarized form")
}
