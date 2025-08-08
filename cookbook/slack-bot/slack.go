package main

import (
	"context"
	"fmt"
	"log"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// SlackService encapsulates all Slack API operations
type SlackService struct {
	client       *slack.Client
	socketClient *socketmode.Client
	botUserID    string
}

// NewSlackService creates a new Slack service
func NewSlackService(botToken, appToken string) (*SlackService, error) {
	// Create Slack clients
	api := slack.New(
		botToken,
		slack.OptionDebug(false),
		slack.OptionAppLevelToken(appToken),
	)

	// Create Socket Mode client
	socketClient := socketmode.New(
		api,
		socketmode.OptionDebug(false),
	)

	// Get bot user ID
	authResp, err := api.AuthTest()
	if err != nil {
		return nil, fmt.Errorf("auth test failed: %w", err)
	}

	service := &SlackService{
		client:       api,
		socketClient: socketClient,
		botUserID:    authResp.UserID,
	}

	log.Printf("Slack service initialized as %s (ID: %s)", authResp.User, service.botUserID)
	return service, nil
}

// RunSocketMode starts the Socket Mode connection
func (s *SlackService) RunSocketMode(ctx context.Context) error {
	return s.socketClient.RunContext(ctx)
}

// GetEvents returns the Socket Mode events channel
func (s *SlackService) GetEvents() <-chan socketmode.Event {
	return s.socketClient.Events
}

// AckEvent acknowledges a Socket Mode event
func (s *SlackService) AckEvent(req socketmode.Request) {
	s.socketClient.Ack(req)
}

// SendMessage sends a message to a Slack channel
func (s *SlackService) SendMessage(channel, text, threadTS string) error {
	options := []slack.MsgOption{
		slack.MsgOptionText(text, false),
	}

	if threadTS != "" {
		options = append(options, slack.MsgOptionTS(threadTS))
	}

	_, _, err := s.client.PostMessage(channel, options...)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

// IsBotMessage checks if a message is from a bot (including self)
func (s *SlackService) IsBotMessage(event *slackevents.MessageEvent) bool {
	return event.User == s.botUserID || event.BotID != ""
}

// IsBotMention checks if an app mention is from a bot
func (s *SlackService) IsBotMention(event *slackevents.AppMentionEvent) bool {
	return event.User == s.botUserID || event.BotID != ""
}

// GetBotUserID returns the bot's user ID
func (s *SlackService) GetBotUserID() string {
	return s.botUserID
}

// GetThreadMessages retrieves messages from a thread
func (s *SlackService) GetThreadMessages(channel, threadTS string) ([]slack.Message, error) {
	params := &slack.GetConversationRepliesParameters{
		ChannelID: channel,
		Timestamp: threadTS,
	}

	messages, _, _, err := s.client.GetConversationReplies(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread messages: %w", err)
	}

	return messages, nil
}

// GetChannelHistory retrieves recent messages from a channel or DM
func (s *SlackService) GetChannelHistory(channel string, limit int) ([]slack.Message, error) {
	params := &slack.GetConversationHistoryParameters{
		ChannelID: channel,
		Limit:     limit,
	}

	resp, err := s.client.GetConversationHistory(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel history: %w", err)
	}

	return resp.Messages, nil
}
