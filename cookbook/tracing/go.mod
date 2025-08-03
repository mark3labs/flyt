module github.com/mark3labs/flyt/cookbook/tracing

go 1.23

toolchain go1.24.4

require (
	github.com/google/uuid v1.6.0
	github.com/henomis/langfuse-go v0.0.3
	github.com/mark3labs/flyt v0.0.0
)

require github.com/henomis/restclientgo v1.2.0 // indirect

replace github.com/mark3labs/flyt => ../../
