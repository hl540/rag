package embedding

import "context"

type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
	Embeds(ctx context.Context, texts []string) ([][]float32, error)
}
