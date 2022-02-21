/*
Copyright © 2022 Rodion Lim <rodion.lim@hotmail.com>

*/
package cmd

import (
	"context"
	"fmt"
	"os"

	crs "github.com/rodionlim/carousell/library/carousell"
	"github.com/rodionlim/carousell/library/log"
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
		logger := log.Ctx(context.Background())
		setVerbosity(cmd)
		opts := getOpts(cmd, args)

		req := crs.NewReq(opts...)
		listings, err := req.Get()
		if err != nil {
			fmt.Println("Something unexpected happened")
			logger.Error(err.Error())
			os.Exit(1)
		}

		// If user specifies that they want a summarized version of the output
		shortFlag, _ := cmd.Flags().GetBool("shorthand")
		if shortFlag {
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

	getCmd.Flags().BoolP("shorthand", "s", false, "Display listings output in summarized form")
}
