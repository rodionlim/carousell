package notifier

import (
	"context"
	"os"

	"github.com/rodionlim/carousell/library/log"
	"github.com/slack-go/slack"
)

type Slack struct {
	client *slack.Client
}

func NewSlacker() *Slack {
	logger := log.Ctx(context.Background())
	token, exists := os.LookupEnv("SLACK_ACCESS_TOKEN")
	if !exists {
		logger.Error("Invalid slack access token. Please set \"SLACK_ACCESS_TOKEN\" variable")
		os.Exit(1)
	}

	return &Slack{
		client: slack.New(token),
	}
}

// Notify acts as a wrapper to send basic slack messages.
// It requires the env SLACK_ACCESS_TOKEN to be set.
func (s *Slack) Notify(channelID string, msg string, attachment *slack.Attachment) {
	logger := log.Ctx(context.Background())
	msgOptions := []slack.MsgOption{
		slack.MsgOptionText(msg, false),
		slack.MsgOptionAsUser(true), // Add this if you want that the bot would post message as a user, otherwise it will send response using the default slackbot
	}
	if attachment != nil {
		msgOptions = append(msgOptions, slack.MsgOptionAttachments(*attachment))
	}
	channelID, timestamp, err := s.client.PostMessage(
		channelID,
		msgOptions...,
	)
	if err != nil {
		logger.Errorf("%s\n", err)
		return
	}
	logger.Infof("Message successfully sent to channel %s at %s\n", channelID, timestamp)
}
