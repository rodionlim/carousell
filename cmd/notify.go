/*
Copyright © 2022 Rodion Lim <rodion.lim@hotmail.com>

*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	crs "github.com/rodionlim/carousell/library/carousell"
	"github.com/rodionlim/carousell/library/log"
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
		setVerbosity(cmd)
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

		logger := log.Ctx(context.Background())
		slacker := notifier.NewSlacker()
		cache := crs.NewCache()
		listings, err := req.Get()
		if err != nil {
			fmt.Println("Something unexpected happened")
			logger.Error(err.Error())
			os.Exit(1)
		}
		logger.Infof("***Recv initial listings*** \n%v", crs.ShortenListings(listings))
		cache.Store(listings)
		logger.Info("Cached initial listings")

		d := time.Minute * time.Duration(interval)
		ticker := time.NewTicker(d)
		for {
			logger.Infof("Waiting for %d mins before next query\n", interval)
			<-ticker.C
			listings, err = req.Get()
			logger.Infof("***Recv listings*** \n%v", crs.ShortenListings(listings))
			if err != nil {
				fmt.Println("Something unexpected happened")
				logger.Error(err.Error())
				os.Exit(1)
			}
			cache.ProcessAndStore(listings, func(listing crs.Listing) error {
				slacker.Notify(slackChannel, listing.Sprint(), nil)
				return nil
			}, true)
		}
	},
}

func init() {
	rootCmd.AddCommand(notifyCmd)

	notifyCmd.Flags().String("slack-channel", "", "Slack channel id to send notifications, e.g. C0341H4MD1P")
	notifyCmd.MarkFlagRequired("slack-channel")

	notifyCmd.Flags().IntP("interval", "i", 10, "Interval in minutes")
}
