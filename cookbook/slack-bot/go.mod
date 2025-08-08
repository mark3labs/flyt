module github.com/mark3labs/flyt/cookbook/slack-bot

go 1.23

toolchain go1.24.5

require (
	github.com/joho/godotenv v1.5.1
	github.com/mark3labs/flyt v0.0.0
	github.com/slack-go/slack v0.12.3
)

require github.com/gorilla/websocket v1.5.0 // indirect

replace github.com/mark3labs/flyt => ../..
