package notifier

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

type Slack struct {
	client  *slack.Client
	channel string // Default channel to post messages to if not specified
}

func NewSlacker() *Slack {
	token, exists := os.LookupEnv("SLACK_ACCESS_TOKEN")
	if !exists {
		fmt.Println("Invalid slack access token. Please set \"SLACK_ACCESS_TOKEN\" variable")
	}

	return &Slack{
		client: slack.New(token),
	}
}

// Notify acts as a wrapper to send basic slack messages.
// It requires the env SLACK_ACCESS_TOKEN to be set.
func (s *Slack) Notify(channelID string, msg string, attachment *slack.Attachment) {
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
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
}
