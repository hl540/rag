package llm

import (
	"context"
	"errors"
	"github.com/ollama/ollama/api"
	"net/http"
	"net/url"
)

type OllamaLLM struct {
	client *api.Client
	model  string
}

func New(base string, model string) (LLM, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	return &OllamaLLM{
		client: api.NewClient(baseURL, http.DefaultClient),
		model:  model,
	}, nil
}

func (l *OllamaLLM) ChatCompletion(ctx context.Context, opts ChatCompletionOptions) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (l *OllamaLLM) CreateEmbedding(ctx context.Context, text string) ([]float32, error) {
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
