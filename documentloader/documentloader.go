package documentloader

import (
	"context"
	"github.com/hl540/rag/textsplitter"
	"github.com/hl540/rag/vectorstore"
)

type DocumentLoader interface {
	Load(ctx context.Context) ([]*vectorstore.Document, error)
	LoadSplit(ctx context.Context, splitter textsplitter.TextSplitter) ([]*vectorstore.Document, error)
}
