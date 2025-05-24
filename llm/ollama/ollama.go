package ollama

import (
	"context"
	"errors"
	"github.com/hl540/rag/llm"
	"github.com/ollama/ollama/api"
	"net/http"
	"net/url"
)

type LLM struct {
	client *api.Client
	model  string
}

func New(base string, model string) (llm.LLM, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	return &LLM{
		client: api.NewClient(baseURL, http.DefaultClient),
		model:  model,
	}, nil
}

func (l *LLM) ChatCompletion(ctx context.Context, opts llm.ChatCompletionOptions) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LLM) CreateEmbedding(ctx context.Context, text string) ([]float32, error) {
	embed, err := l.client.Embed(ctx, &api.EmbedRequest{
		Model: l.model,
		Input: text,
	})
	if err != nil {
		return nil, err
	}
	if len(embed.Embeddings) == 0 {
		return nil, errors.New("no embeddings")
	}
	return embed.Embeddings[0], nil
}
