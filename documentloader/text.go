package documentloader

import (
	"bufio"
	"context"
	"github.com/google/uuid"
	"github.com/hl540/rag/textsplitter"
	"github.com/hl540/rag/vectorstore"
	"io"
	"strings"
)

type TextLoader struct {
	r io.Reader
}

func New(read io.Reader) DocumentLoader {
	return &TextLoader{
		r: read,
	}
}

func (l *TextLoader) Load(ctx context.Context) ([]*vectorstore.Document, error) {
	docs := make([]*vectorstore.Document, 0)
	scanner := bufio.NewScanner(l.r)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		docs = append(docs, &vectorstore.Document{
			Id:   uuid.NewString(),
			Text: line,
			Metadata: map[string]any{
				vectorstore.ContentKey: line,
			},
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return docs, nil
}

func (l *TextLoader) LoadSplit(ctx context.Context, splitter textsplitter.TextSplitter) ([]*vectorstore.Document, error) {
	all, err := io.ReadAll(l.r)
	if err != nil {
		return nil, err
	}
	docs := make([]*vectorstore.Document, 0)
	for _, text := range splitter.SplitText(string(all)) {
		if text == "" {
			continue
		}
		docs = append(docs, &vectorstore.Document{
			Id:   uuid.NewString(),
			Text: text,
			Metadata: map[string]any{
				vectorstore.ContentKey: text,
			},
		})
	}
	return docs, nil
}
