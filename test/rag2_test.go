package test

import (
	"context"
	"github.com/hl540/rag/documentloader"
	"github.com/hl540/rag/embedding/ollama"
	"github.com/hl540/rag/vectorstore/qdrant"
	"github.com/ollama/ollama/api"
	"log"
	"os"
	"testing"
)

func TestEmbed(t *testing.T) {
	llm, err := api.ClientFromEnvironment()
	if err != nil {
		t.Fatalf("ollama 连接失败：%s", err.Error())
	}

	ctx := context.Background()
	vectorStore, err := qdrant.New(
		qdrant.WithHost("106.55.106.158"),
		qdrant.WithPort(6334),
		qdrant.WithEmbedder(ollama.New(llm, "quentinz/bge-base-zh-v1.5:latest")),
	)
	//vectorStore := memory.New(ollama.New(llm, "quentinz/bge-base-zh-v1.5:latest"))

	file, err := os.Open("../szbf.txt")
	if err != nil {
		t.Fatalf("文件读取失败：%s", err.Error())
	}
	loader := documentloader.New(file)
	//docs, err := loader.LoadSplit(ctx, textsplitter.NewSentenceSplitter(500, 30, false))
	docs, err := loader.Load(ctx)
	if err != nil {
		log.Fatalf("文件加载失败： %s", err.Error())
	}
	err = vectorStore.AddDocuments(ctx, "szbf", docs)
	if err != nil {
		log.Fatalf("向量化存储失败：%s", err.Error())
	}
	t.Log("Success")
	search, err := vectorStore.SimilaritySearch(ctx, "szbf", "孙子兵法“军形篇”讲了什么", 5)
	if err != nil {
		t.Fatalf("vector 查询失败: %v", err)
	}
	for _, doc := range search {
		t.Logf("ID:%s, Score:%f, Metadata:%+v", doc.ID, doc.Score, doc.Metadata)
	}
}
