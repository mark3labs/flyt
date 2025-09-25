package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/mark3labs/flyt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get required tokens
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")
	openAIKey := os.Getenv("OPENAI_API_KEY")

	if botToken == "" || appToken == "" || openAIKey == "" {
		log.Fatal("Missing required environment variables: SLACK_BOT_TOKEN, SLACK_APP_TOKEN, or OPENAI_API_KEY")
	}

	// Create Slack service
	slackService, err := NewSlackService(botToken, appToken)
	if err != nil {
		log.Fatalf("Failed to create Slack service: %v", err)
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		cancel()
	}()

	// Start the bot
	bot := &SlackBot{
		slack:       slackService,
		openAIKey:   openAIKey,
		llmServices: make(map[string]*LLMService),
	}

	log.Println("🤖 Slack Bot with Flyt starting...")
	if err := bot.Start(ctx); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}

type SlackBot struct {
	slack       *SlackService
	openAIKey   string
	llmServices map[string]*LLMService
	mu          sync.RWMutex
}

func (b *SlackBot) Start(ctx context.Context) error {
	// Start event handler
	go b.handleEvents(ctx)

	// Start cleanup routine for old conversations
	go b.cleanupOldConversations(ctx)

	// Run Socket Mode client
	log.Println("Connected to Slack with Socket Mode")
	return b.slack.RunSocketMode(ctx)
}
func (b *SlackBot) cleanupOldConversations(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.mu.Lock()
			// Simple memory management: clear all conversations if we have too many
			// In production, implement proper LRU cache or track last access time
			if len(b.llmServices) > 100 {
				log.Printf("Memory cleanup: clearing %d conversations", len(b.llmServices))
				b.llmServices = make(map[string]*LLMService)
			}
			b.mu.Unlock()
		}
	}
}
func (b *SlackBot) handleEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case evt := <-b.slack.GetEvents():
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				log.Println("Connecting to Slack...")
			case socketmode.EventTypeConnected:
				log.Println("Connected to Slack")
			case socketmode.EventTypeConnectionError:
				log.Printf("Connection error: %v", evt.Data)
			case socketmode.EventTypeEventsAPI:
				b.handleEventAPI(ctx, evt)
			}
		}
	}
}

func (b *SlackBot) handleEventAPI(ctx context.Context, evt socketmode.Event) {
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		log.Printf("Ignored non-EventsAPI event: %v", evt.Type)
		return
	}

	// Acknowledge the event
	b.slack.AckEvent(*evt.Request)

	switch eventsAPIEvent.Type {
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			// Skip bot's own messages
			if b.slack.IsBotMessage(ev) {
				return
			}
			b.handleMessage(ctx, ev)
		case *slackevents.AppMentionEvent:
			// Skip bot's own mentions
			if b.slack.IsBotMention(ev) {
				return
			}
			b.handleMention(ctx, ev)
		}
	}
}

func (b *SlackBot) handleMessage(ctx context.Context, event *slackevents.MessageEvent) {
	// Only process direct messages (not in a channel)
	if !b.isDirectMessage(event.Channel) {
		// In channels, we only respond to mentions (handled by handleMention)
		return
	}

	log.Printf("DM from %s: %s", event.User, event.Text)

	// For DMs, don't use threads - respond directly in the conversation
	// Pass empty string for threadTS to indicate no threading
	b.processWithFlyt(ctx, event.Text, event.Channel, "")
}

func (b *SlackBot) handleMention(ctx context.Context, event *slackevents.AppMentionEvent) {
	log.Printf("Mention from %s in channel %s: %s", event.User, event.Channel, event.Text)

	// For mentions in channels, always reply in thread
	threadTS := event.ThreadTimeStamp
	if threadTS == "" {
		// Start a new thread with the mention message as root
		threadTS = event.TimeStamp
	}

	// Process mention through Flyt workflow
	b.processWithFlyt(ctx, event.Text, event.Channel, threadTS)
}

func (b *SlackBot) isDirectMessage(channel string) bool {
	// Direct message channels start with 'D'
	// Group DMs start with 'G'
	return len(channel) > 0 && (channel[0] == 'D' || channel[0] == 'G')
}
func (b *SlackBot) getLLMService(channel, threadTS string) *LLMService {
	// Create a key for this conversation context
	key := channel
	if threadTS != "" {
		key = channel + ":" + threadTS
	}

	b.mu.RLock()
	service, exists := b.llmServices[key]
	b.mu.RUnlock()

	if !exists {
		b.mu.Lock()
		service = NewLLMService(b.openAIKey)
		b.llmServices[key] = service
		b.mu.Unlock()
		log.Printf("Created new LLM service for conversation: %s", key)
	}

	return service
}

func (b *SlackBot) processWithFlyt(ctx context.Context, message, channel, threadTS string) {
	// Get or create LLM service for this thread
	llmService := b.getLLMService(channel, threadTS)

	// Fetch conversation history for context
	history := b.fetchConversationHistory(channel, threadTS)

	// Create shared store
	shared := flyt.NewSharedStore()
	shared.Set("message", message)
	shared.Set("channel", channel)
	shared.Set("thread_ts", threadTS)
	shared.Set("history", history)
	// Create workflow with injected LLM service
	flow := b.createWorkflow(llmService)

	// Run workflow
	if err := flow.Run(ctx, shared); err != nil {
		log.Printf("Workflow error: %v", err)
		b.sendMessage(channel, "Sorry, I encountered an error processing your request.", threadTS)
		return
	}

	// Get response from shared store
	if response := shared.GetString("response"); response != "" {
		b.sendMessage(channel, response, threadTS)
	}
}

func (b *SlackBot) fetchConversationHistory(channel, threadTS string) []map[string]string {
	var history []map[string]string
	var messages []slack.Message
	var err error

	if threadTS != "" {
		// For threads in channels, fetch thread messages
		messages, err = b.slack.GetThreadMessages(channel, threadTS)
		if err != nil {
			log.Printf("Failed to fetch thread history: %v", err)
			return history
		}
	} else if b.isDirectMessage(channel) {
		// For DMs, fetch recent channel history
		messages, err = b.slack.GetChannelHistory(channel, 20)
		if err != nil {
			log.Printf("Failed to fetch DM history: %v", err)
			return history
		}
	}

	// Convert to simplified format, excluding bot's own messages
	botID := b.slack.GetBotUserID()
	for _, msg := range messages {
		// Skip bot's own messages and empty messages
		if msg.User == botID || msg.BotID != "" || msg.Text == "" {
			continue
		}

		history = append(history, map[string]string{
			"user":      msg.User,
			"text":      msg.Text,
			"timestamp": msg.Timestamp,
		})
	}

	// Limit history to last 10 messages for context
	if len(history) > 10 {
		history = history[len(history)-10:]
	}

	return history
}

func (b *SlackBot) createWorkflow(llmService *LLMService) *flyt.Flow {
	// Create nodes with injected dependencies
	parseNode := NewParseMessageNode(b.slack)
	llmNode := NewLLMNode(llmService)
	toolNode := NewToolExecutorNode()
	formatNode := NewFormatResponseNode()

	// Create flow
	flow := flyt.NewFlow(parseNode).
		Connect(parseNode, flyt.DefaultAction, llmNode).
		Connect(llmNode, "tool_call", toolNode).
		Connect(llmNode, "response", formatNode).
		Connect(toolNode, flyt.DefaultAction, llmNode).
		Connect(formatNode, flyt.DefaultAction, nil)

	return flow
}

func (b *SlackBot) sendMessage(channel, text, threadTS string) {
	if err := b.slack.SendMessage(channel, text, threadTS); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
