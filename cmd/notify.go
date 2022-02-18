/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"time"

	crs "github.com/rodionlim/carousell/library/carousell"
	"github.com/rodionlim/carousell/library/notifier"
	"github.com/spf13/cobra"
)

// notifyCmd represents the notify command
var notifyCmd = &cobra.Command{
	Use:   "notify [search terms...]",
	Short: "Notify user on new carousell listings",
	Long: `
Notify user on new carousell listings. Arguments supplied to this command
are search terms in Carousell, and at least one search term should be provided.

If using an app's oauth access token, ensure that the app has been invited to the channel.

Flags can be used to modify the search behaviour, e.g. specifying the 
-r recent flag, to include only recent listings etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := getOpts(cmd, args)
		req := crs.NewReq(opts...)
		interval, _ := cmd.Flags().GetInt("interval")
		slackChannel, _ := cmd.Flags().GetString("slack-channel")

		// Initialization log
		fmt.Printf(`
***
Setting up slack notifications with parameters:
Search Terms: %s
Interval: %d
Slack Channel: %s
***

`, req.GetSearchTerm(), interval, slackChannel)

		slacker := notifier.NewSlacker()
		cache := crs.NewCache()
		listings, err := req.Get()
		if err != nil {
			fmt.Println("Something unexpected happened")
			os.Exit(1)
		}
		cache.Store(listings)

		d := time.Minute * time.Duration(interval)
		ticker := time.NewTicker(d)
		for {
			fmt.Printf("Waiting for %d mins before next query\n", interval)
			<-ticker.C
			listings, err = req.Get()
			if err != nil {
				fmt.Println("Something unexpected happened")
				os.Exit(1)
			}
			cache.ProcessAndStore(listings, func(listing crs.Listing) error {
				slacker.Notify(slackChannel, fmt.Sprintf("%s - $%f", listing.Title, listing.Price), nil)
				return nil
			})
		}
	},
}

func init() {
	rootCmd.AddCommand(notifyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// notifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// notifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	notifyCmd.Flags().String("slack-channel", "", "Slack channel id to send notifications, e.g. C0341H4MD1P")
	notifyCmd.MarkFlagRequired("slack-channel")

	notifyCmd.Flags().IntP("interval", "i", 10, "Interval in minutes")
}
