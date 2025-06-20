package documentloader

import (
	"github.com/hl540/rag/textsplitter"
	"github.com/hl540/rag/vectorstore"
)

type DocumentLoader interface {
	Load() ([]*vectorstore.Document, error)
	LoadSplit(splitter textsplitter.TextSplitter) ([]*vectorstore.Document, error)
}
