# Slack Bot with OpenAI GPT-4.1 and Flyt

A sophisticated Slack bot built with Flyt workflow framework that integrates OpenAI GPT-4.1 with function calling capabilities. The bot can perform calculations and share Chuck Norris facts while maintaining conversation context.

## Features

- **OpenAI GPT-4.1 Integration**: Uses GPT-4.1 model with function calling API
- **Tool Support**: 
  - Calculator for mathematical expressions
  - Chuck Norris fact generator for entertainment
- **Flyt Workflow**: Leverages Flyt's node-based architecture for clean separation of concerns
- **Socket Mode**: Real-time message handling without webhooks
- **Conversation Memory**: Maintains context across messages
- **Thread Support**: Responds in threads when appropriate

## Architecture

The bot uses Flyt's workflow pattern with the following nodes:

1. **ParseMessageNode**: Cleans and prepares incoming messages
2. **LLMNode**: Processes messages through OpenAI with function calling (with injected LLM service)
3. **ToolExecutorNode**: Executes requested tools (calculator, Chuck Norris facts)
4. **FormatResponseNode**: Formats responses for Slack

```
User Message → Parse → LLM → Tool Execution (if needed) → Format → Response
                         ↑                    ↓
                         └────────────────────┘
```

### Key Design Features

- **Service Abstraction**: Both Slack and LLM operations are encapsulated in service objects
- **Dependency Injection**: Services are injected into nodes that need them, not passed through SharedStore
- **Per-Thread Conversations**: Each Slack thread maintains its own conversation history
- **Memory Management**: Automatic cleanup of old conversations after 30 minutes of inactivity
- **Clean Separation**: Configuration (API keys) and external services separate from workflow data
- **Testability**: Services can be easily mocked for unit testing

## Prerequisites

- Go 1.21 or higher
- Slack workspace with admin access
- OpenAI API key
- Slack Bot and App tokens

## Slack App Setup

### Quick Setup with App Manifest (Recommended)

1. Go to [api.slack.com/apps](https://api.slack.com/apps)
2. Click "Create New App" → "From an app manifest"
3. Select your workspace
4. Choose "YAML" format and paste the following manifest:

```yaml
display_information:
  name: Flyt AI Bot
  description: An AI-powered Slack bot using OpenAI GPT-4.1 with calculator and Chuck Norris facts
  background_color: "#2c2d30"
  long_description: This bot integrates OpenAI's GPT-4.1 with function calling capabilities to provide intelligent responses, perform calculations, and share Chuck Norris facts. Built with the Flyt workflow framework for clean, maintainable code architecture.
features:
  app_home:
    home_tab_enabled: false
    messages_tab_enabled: true
    messages_tab_read_only_enabled: false
  bot_user:
    display_name: Flyt AI Bot
    always_online: true
oauth_config:
  scopes:
    bot:
      - channels:history
      - channels:read
      - chat:write
      - groups:history
      - groups:read
      - im:history
      - im:read
      - im:write
      - mpim:history
      - mpim:read
      - mpim:write
      - users:read
      - app_mentions:read
settings:
  event_subscriptions:
    bot_events:
      - app_mention
      - message.channels
      - message.groups
      - message.im
      - message.mpim
  interactivity:
    is_enabled: false
  org_deploy_enabled: false
  socket_mode_enabled: true
  token_rotation_enabled: false
```

5. Click "Next" and review the configuration
6. Click "Create"

### After App Creation

1. **Generate App-Level Token**:
   - Go to "Basic Information" → "App-Level Tokens"
   - Click "Generate Token and Scopes"
   - Name it (e.g., "Socket Mode Token")
   - Add the `connections:write` scope
   - Click "Generate"
   - Save this token as `SLACK_APP_TOKEN`

2. **Install App to Workspace**:
   - Go to "Install App" in the sidebar
   - Click "Install to Workspace"
   - Authorize the app
   - Copy the Bot User OAuth Token as `SLACK_BOT_TOKEN`

### Manual Setup (Alternative)

If you prefer to set up manually or need to modify an existing app:

1. **Create a Slack App**:
   - Go to [api.slack.com/apps](https://api.slack.com/apps)
   - Click "Create New App" → "From scratch"
   - Name your app and select workspace

2. **Configure OAuth & Permissions**:
   Add these Bot Token Scopes:
   - `channels:history`, `channels:read`, `chat:write`
   - `groups:history`, `groups:read`
   - `im:history`, `im:read`, `im:write`
   - `mpim:history`, `mpim:read`, `mpim:write`
   - `users:read`, `app_mentions:read`

3. **Enable Socket Mode**:
   - Go to "Socket Mode" in the sidebar
   - Enable Socket Mode
   - Generate an App-Level Token with `connections:write` scope

4. **Subscribe to Events**:
   - Enable Event Subscriptions
   - Add bot events: `app_mention`, `message.channels`, `message.groups`, `message.im`, `message.mpim`

5. **Install App**:
   - Install to workspace and copy tokens

## Installation

1. Clone the repository:
```bash
git clone https://github.com/mark3labs/flyt
cd flyt/cookbook/slack-bot
```

2. Install dependencies:
```bash
go mod tidy
```

3. Create a `.env` file:
```env
# Required
SLACK_BOT_TOKEN=xoxb-your-bot-token
SLACK_APP_TOKEN=xapp-your-app-token
OPENAI_API_KEY=sk-your-openai-key

# Optional
LOG_LEVEL=info
```

## Usage

### Running the Bot

```bash
go run .
```

The bot will connect to Slack and start listening for messages.

### Interacting with the Bot

1. **Direct Message**: Send a DM to the bot
2. **Channel Mention**: @YourBot in any channel
3. **Channel Message**: Any message in channels where the bot is present

### Example Interactions

**Calculator:**
```
User: Can you calculate 25 * 4 + sqrt(16)?
Bot: I'll help you calculate that expression.
     Result: 104
```

**Chuck Norris Facts:**
```
User: Tell me a Chuck Norris fact about programming
Bot: Here's a Chuck Norris fact for you:
     Chuck Norris doesn't use web standards. The web conforms to Chuck Norris.
```

**Combined:**
```
User: What's 2^10 and also give me a Chuck Norris fact
Bot: I'll calculate 2^10 and get you a Chuck Norris fact.
     
     The calculation 2^10 equals 1024.
     
     And here's your Chuck Norris fact:
     Chuck Norris can divide by zero.
```

## How It Works

### Function Calling Flow

1. User sends a message to the bot
2. Bot retrieves or creates an LLM service for the specific thread/channel
3. The message is processed through the Flyt workflow with injected LLM service
4. OpenAI GPT-4.1 analyzes the message and determines if tools are needed
5. If tools are requested, the bot executes them and sends results back to GPT-4.1
6. GPT-4.1 formulates a final response using the tool results
7. The response is sent back to the user in Slack
8. Conversation history is maintained per thread for context continuity

### Available Tools

#### Calculator
- Evaluates mathematical expressions
- Supports basic operations: +, -, *, /
- Functions: sqrt(), pow()
- Example: `"2 + 2"`, `"sqrt(16)"`, `"pow(2, 8)"`

#### Chuck Norris Facts
- Returns random Chuck Norris facts
- Optional categories: dev, movie, food, sport
- Falls back to hardcoded facts if API is unavailable

## Code Structure

```
slack-bot/
├── main.go       # Entry point and bot orchestration
├── slack.go      # Slack service abstraction
├── llm.go        # LLM service with OpenAI GPT-4.1 integration
├── tools.go      # Tool implementations (calculator, Chuck Norris)
├── nodes.go      # Flyt node implementations
├── manifest.yml  # Slack app manifest for easy setup
├── go.mod        # Go module definition
├── go.sum        # Dependency checksums
├── .env.example  # Environment variable template
└── README.md     # This file
```

### Key Components

**main.go**
- Bot orchestration and lifecycle management
- Per-thread LLM service management
- Workflow creation with dependency injection
- Memory cleanup routines

**slack.go**
- SlackService: Encapsulates all Slack API operations
- Message sending, reactions, thread management
- User and channel information retrieval
- Socket Mode event handling abstraction

**llm.go**
- LLMService: Manages OpenAI client and conversations
- Function calling definitions for tools
- Conversation history management
- Request/response handling with GPT-4.1

**tools.go**
- Calculator expression evaluator
- Chuck Norris API integration
- Tool execution dispatcher

**nodes.go**
- ParseMessageNode: Message preprocessing with optional Slack context
- LLMNode: AI processing with injected LLM service
- ToolExecutorNode: Tool execution
- FormatResponseNode: Response formatting with Slack service integration

## Configuration

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `SLACK_BOT_TOKEN` | Yes | Bot User OAuth Token |
| `SLACK_APP_TOKEN` | Yes | App-Level Token for Socket Mode |
| `OPENAI_API_KEY` | Yes | OpenAI API key |
| `LOG_LEVEL` | No | Logging level (debug, info, warn, error) |

### OpenAI Settings

The bot uses GPT-4 with these settings:
- Model: `gpt-4.1`
- Temperature: 0.7
- Tool Choice: `auto`
- Max conversation history: 20 messages

## Error Handling

The bot includes comprehensive error handling:
- Graceful shutdown on SIGINT/SIGTERM
- Automatic reconnection on connection loss
- Fallback responses for API failures
- Tool execution error recovery
- Context cancellation support

## Development

### Testing Tools Locally

```go
// Test calculator
executor := NewToolExecutor()
result, err := executor.ExecuteTool("calculator", `{"expression": "2 + 2"}`)

// Test Chuck Norris facts
result, err := executor.ExecuteTool("chuck_norris_fact", `{"category": "dev"}`)
```

### Adding New Tools

1. Define the tool in `getToolDefinitions()` in `llm.go`
2. Implement the tool execution in `tools.go`
3. Add the case in `ExecuteTool()` switch statement

Example:
```go
// In llm.go
{
    Type: "function",
    Function: FunctionDef{
        Name:        "weather",
        Description: "Get weather information",
        Parameters: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "location": map[string]interface{}{
                    "type": "string",
                    "description": "City name",
                },
            },
            "required": []string{"location"},
        },
    },
}

// In tools.go
case "weather":
    return te.executeWeather(arguments)
```

## Troubleshooting

### Bot not responding
- Check Socket Mode is enabled in Slack app settings
- Verify bot is invited to the channel
- Check logs for connection errors
- Ensure tokens are correct in `.env`

### OpenAI errors
- Verify API key is valid
- Check API rate limits
- Ensure you have GPT-4 access
- Review error logs for specific issues

### Tool execution failures
- Calculator: Check expression syntax
- Chuck Norris: API may be down, fallback will be used
- Check logs for detailed error messages

## Security Considerations

- Never commit `.env` file or tokens
- Use environment variables for sensitive data
- Rotate tokens regularly
- Limit bot permissions to necessary scopes
- Monitor API usage and costs

## Performance

- Concurrent message processing via Flyt
- Connection pooling for HTTP requests
- Conversation history trimming (max 20 messages)
- Efficient tool execution with timeouts
- Graceful error handling without retries

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - See the main Flyt repository for details.

## Acknowledgments

- [Flyt](https://github.com/mark3labs/flyt) - Workflow framework
- [slack-go](https://github.com/slack-go/slack) - Slack SDK
- [OpenAI](https://openai.com) - GPT-4 API
- [Chuck Norris API](https://api.chucknorris.io) - Chuck Norris facts