package llm

import "context"

type ChatCompletionOptions struct {
	Messages []Message
	// ...
}

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)

type Message struct {
	Role    string // "user", "assistant", "system"
	Content string
}

type LLM interface {
	ChatCompletion(ctx context.Context, opts ChatCompletionOptions) (string, error)
	CreateEmbedding(ctx context.Context, text string) ([]float32, error)
}
