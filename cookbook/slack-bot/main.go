package main

import (
	"context"
	"fmt"
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

	// Create Slack clients
	api := slack.New(
		botToken,
		slack.OptionDebug(false),
		slack.OptionAppLevelToken(appToken),
	)

	// Create Socket Mode client
	client := socketmode.New(
		api,
		socketmode.OptionDebug(false),
	)

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
		client:       api,
		socketClient: client,
		openAIKey:    openAIKey,
		llmServices:  make(map[string]*LLMService),
	}

	log.Println("🤖 Slack Bot with Flyt starting...")
	if err := bot.Start(ctx); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}

type SlackBot struct {
	client       *slack.Client
	socketClient *socketmode.Client
	openAIKey    string
	botUserID    string
	llmServices  map[string]*LLMService
	mu           sync.RWMutex
}

func (b *SlackBot) Start(ctx context.Context) error {
	// Get bot user ID
	authResp, err := b.client.AuthTest()
	if err != nil {
		return fmt.Errorf("auth test failed: %w", err)
	}
	b.botUserID = authResp.UserID
	log.Printf("Bot authenticated as %s (ID: %s)", authResp.User, b.botUserID)

	// Start event handler
	go b.handleEvents(ctx)

	// Start cleanup routine for old conversations
	go b.cleanupOldConversations(ctx, 30*time.Minute)

	// Run Socket Mode client
	log.Println("Connected to Slack with Socket Mode")
	return b.socketClient.RunContext(ctx)
}

func (b *SlackBot) cleanupOldConversations(ctx context.Context, maxAge time.Duration) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.mu.Lock()
			// In a production system, you'd track last access time
			// For now, we'll just clear if the map gets too large
			if len(b.llmServices) > 100 {
				log.Printf("Clearing %d old conversations", len(b.llmServices))
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
		case evt := <-b.socketClient.Events:
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
	b.socketClient.Ack(*evt.Request)

	switch eventsAPIEvent.Type {
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			// Skip bot's own messages
			if ev.User == b.botUserID || ev.BotID != "" {
				return
			}
			b.handleMessage(ctx, ev)
		case *slackevents.AppMentionEvent:
			b.handleMention(ctx, ev)
		}
	}
}

func (b *SlackBot) handleMessage(ctx context.Context, event *slackevents.MessageEvent) {
	log.Printf("Message from %s in channel %s: %s", event.User, event.Channel, event.Text)

	// Process message through Flyt workflow
	b.processWithFlyt(ctx, event.Text, event.Channel, event.ThreadTimeStamp)
}

func (b *SlackBot) handleMention(ctx context.Context, event *slackevents.AppMentionEvent) {
	log.Printf("Mention from %s in channel %s: %s", event.User, event.Channel, event.Text)

	// Process mention through Flyt workflow
	b.processWithFlyt(ctx, event.Text, event.Channel, event.ThreadTimeStamp)
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

	// Create shared store
	shared := flyt.NewSharedStore()
	shared.Set("message", message)
	shared.Set("channel", channel)
	shared.Set("thread_ts", threadTS)

	// Create workflow with injected LLM service
	flow := b.createWorkflow(llmService)

	// Run workflow
	if err := flow.Run(ctx, shared); err != nil {
		log.Printf("Workflow error: %v", err)
		b.sendMessage(channel, "Sorry, I encountered an error processing your request.", threadTS)
		return
	}

	// Get response from shared store
	if response, ok := shared.Get("response"); ok {
		if responseStr, ok := response.(string); ok {
			b.sendMessage(channel, responseStr, threadTS)
		}
	}
}

func (b *SlackBot) createWorkflow(llmService *LLMService) *flyt.Flow {
	// Create nodes with injected dependencies
	parseNode := &ParseMessageNode{BaseNode: flyt.NewBaseNode()}
	llmNode := &LLMNode{
		BaseNode: flyt.NewBaseNode(),
		llm:      llmService,
	}
	toolNode := &ToolExecutorNode{BaseNode: flyt.NewBaseNode()}
	formatNode := &FormatResponseNode{BaseNode: flyt.NewBaseNode()}

	// Create flow
	flow := flyt.NewFlow(parseNode)
	flow.Connect(parseNode, flyt.DefaultAction, llmNode)
	flow.Connect(llmNode, "tool_call", toolNode)
	flow.Connect(llmNode, "response", formatNode)
	flow.Connect(toolNode, flyt.DefaultAction, llmNode)
	flow.Connect(formatNode, flyt.DefaultAction, nil)

	return flow
}

func (b *SlackBot) sendMessage(channel, text, threadTS string) {
	options := []slack.MsgOption{
		slack.MsgOptionText(text, false),
	}

	if threadTS != "" {
		options = append(options, slack.MsgOptionTS(threadTS))
	}

	_, _, err := b.client.PostMessage(channel, options...)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
