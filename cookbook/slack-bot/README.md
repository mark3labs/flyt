# Slack Bot with OpenAI GPT-4.1 and Flyt

A Slack bot built with Flyt that integrates OpenAI GPT-4.1 with function calling to perform calculations and share Chuck Norris facts while maintaining conversation context.

## Features

- **OpenAI GPT-4.1**: Function calling with calculator and Chuck Norris facts tools
- **Smart Threading**: Responds to @mentions in threads (channels) or directly (DMs)
- **Context Awareness**: Maintains conversation history per thread
- **Socket Mode**: Real-time handling without webhooks
- **Auto-cleanup**: Removes inactive conversations after 30 minutes

## Architecture

```mermaid
graph LR
    A[User Message] --> B[ParseMessageNode]
    B --> C[LLMNode]
    C --> D{Tools Needed?}
    D -->|Yes| E[ToolExecutorNode]
    D -->|No| F[FormatResponseNode]
    E --> C
    F --> G[Slack Response]
```

**Nodes:**
- `ParseMessageNode`: Cleans incoming messages
- `LLMNode`: Processes with OpenAI (injected service)
- `ToolExecutorNode`: Executes tools
- `FormatResponseNode`: Formats for Slack

## Prerequisites

- Go 1.21+
- Slack workspace admin access
- OpenAI API key
- Slack Bot/App tokens

## Slack App Setup

1. Go to [api.slack.com/apps](https://api.slack.com/apps)
2. Click "Create New App" → "From an app manifest"
3. Select workspace, choose YAML, paste manifest from `manifest.yml`:

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

4. Create app
5. Generate App-Level Token (Basic Information → App-Level Tokens → Add `connections:write` scope)
6. Install to workspace (Install App → copy Bot User OAuth Token)

## Installation

```bash
git clone https://github.com/mark3labs/flyt
cd flyt/cookbook/slack-bot
go mod tidy
```

Create `.env`:
```env
SLACK_BOT_TOKEN=xoxb-your-bot-token
SLACK_APP_TOKEN=xapp-your-app-token
OPENAI_API_KEY=sk-your-openai-key
LOG_LEVEL=info  # optional
```

## Usage

```bash
go run .
```

### Interactions

- **Channels**: @mention bot → replies in thread
- **DMs**: Direct conversation with context
- **Tools**: Ask for calculations or Chuck Norris facts

Example:
```
User: @Bot calculate sqrt(144) + 8
Bot: [In thread] Result: 20

User: Tell me a Chuck Norris fact
Bot: Chuck Norris can divide by zero.
```

## Available Tools

- **Calculator**: Math expressions with +, -, *, /, sqrt(), pow()
- **Chuck Norris Facts**: Random facts with optional categories (dev, movie, food, sport)

## Project Structure

- `main.go`: Bot orchestration, per-thread LLM management
- `slack.go`: Slack API operations and Socket Mode handling
- `llm.go`: OpenAI GPT-4.1 client with function calling
- `tools.go`: Calculator and Chuck Norris fact implementations
- `nodes.go`: Flyt workflow nodes
- `manifest.yml`: Slack app configuration

## Configuration

**Required Environment Variables:**
- `SLACK_BOT_TOKEN`: Bot User OAuth Token
- `SLACK_APP_TOKEN`: App-Level Token for Socket Mode
- `OPENAI_API_KEY`: OpenAI API key

**OpenAI Settings:**
- Model: `gpt-4.1`, Temperature: 0.7, Max history: 20 messages

## Adding New Tools

1. Define in `getToolDefinitions()` (llm.go)
2. Implement in `ExecuteTool()` (tools.go)

Example:
```go
// llm.go
{
    Type: "function",
    Function: FunctionDef{
        Name: "weather",
        Description: "Get weather",
        Parameters: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "location": map[string]interface{}{"type": "string"},
            },
            "required": []string{"location"},
        },
    },
}

// tools.go
case "weather":
    return te.executeWeather(arguments)
```

## Troubleshooting

- **Bot not responding**: Check Socket Mode enabled, bot invited to channel, tokens valid
- **OpenAI errors**: Verify API key, check rate limits, ensure GPT-4 access
- **Tool failures**: Check expression syntax (calculator), fallback used if Chuck Norris API down

## License

MIT License - See the main Flyt repository for details.