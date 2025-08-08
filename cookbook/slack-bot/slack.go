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

// SendBlocks sends a message with Block Kit blocks
func (s *SlackService) SendBlocks(channel string, blocks []slack.Block, threadTS string) error {
	options := []slack.MsgOption{
		slack.MsgOptionBlocks(blocks...),
	}

	if threadTS != "" {
		options = append(options, slack.MsgOptionTS(threadTS))
	}

	_, _, err := s.client.PostMessage(channel, options...)
	if err != nil {
		return fmt.Errorf("failed to send blocks: %w", err)
	}
	return nil
}

// AddReaction adds an emoji reaction to a message
func (s *SlackService) AddReaction(channel, timestamp, emoji string) error {
	err := s.client.AddReaction(emoji, slack.ItemRef{
		Channel:   channel,
		Timestamp: timestamp,
	})
	if err != nil {
		return fmt.Errorf("failed to add reaction: %w", err)
	}
	return nil
}

// GetUserInfo retrieves information about a user
func (s *SlackService) GetUserInfo(userID string) (*slack.User, error) {
	user, err := s.client.GetUserInfo(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	return user, nil
}

// GetChannelInfo retrieves information about a channel
func (s *SlackService) GetChannelInfo(channelID string) (*slack.Channel, error) {
	channel, err := s.client.GetConversationInfo(&slack.GetConversationInfoInput{
		ChannelID: channelID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get channel info: %w", err)
	}
	return channel, nil
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

// UpdateMessage updates an existing message
func (s *SlackService) UpdateMessage(channel, timestamp, text string) error {
	_, _, _, err := s.client.UpdateMessage(
		channel,
		timestamp,
		slack.MsgOptionText(text, false),
	)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}
	return nil
}

// DeleteMessage deletes a message
func (s *SlackService) DeleteMessage(channel, timestamp string) error {
	_, _, err := s.client.DeleteMessage(channel, timestamp)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}

// SetTyping shows typing indicator in a channel
func (s *SlackService) SetTyping(channel string) error {
	// Note: Slack doesn't have a direct API for this in the SDK
	// This is a placeholder for future implementation
	return nil
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

// EnrichMessageContext adds user and channel information to the message context
func (s *SlackService) EnrichMessageContext(userID, channelID string) map[string]string {
	context := make(map[string]string)

	// Get user info
	if user, err := s.GetUserInfo(userID); err == nil {
		context["user_name"] = user.Name
		context["user_real_name"] = user.RealName
		context["user_tz"] = user.TZ
	}

	// Get channel info
	if channel, err := s.GetChannelInfo(channelID); err == nil {
		context["channel_name"] = channel.Name
		context["channel_topic"] = channel.Topic.Value
		context["channel_purpose"] = channel.Purpose.Value
	}

	return context
}

// JoinChannel joins a channel
func (s *SlackService) JoinChannel(channelID string) error {
	_, _, _, err := s.client.JoinConversation(channelID)
	if err != nil {
		return fmt.Errorf("failed to join channel: %w", err)
	}
	return nil
}

// LeaveChannel leaves a channel
func (s *SlackService) LeaveChannel(channelID string) error {
	_, err := s.client.LeaveConversation(channelID)
	if err != nil {
		return fmt.Errorf("failed to leave channel: %w", err)
	}
	return nil
}
