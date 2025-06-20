package test

import (
	"context"
	"github.com/hl540/rag/embedding"
	"github.com/hl540/rag/vectorstore"
	"github.com/ollama/ollama/api"
	"testing"
)

func TestSGYYSearch(t *testing.T) {
	// 1. 初始化 Ollama 客户端
	llm, err := api.ClientFromEnvironment()
	if err != nil {
		t.Fatalf("Failed to create Ollama client: %v", err)
	}

	// 2. 创建文本分割器 - 使用较小的块大小以保持上下文完整性
	//splitter, err := textsplitter.NewSentenceSplitter(200, 20, true)
	//splitter, err := textsplitter.NewRecursiveCharacterTextSplitterWithDefaults(200, 20)
	//if err != nil {
	//	t.Fatalf("Failed to create text splitter: %v", err)
	//}

	// 3. 创建向量嵌入器 - 使用 bge-base-zh-v1.5 模型
	embedder := embedding.NewOllamaEmbedder(llm, "quentinz/bge-base-zh-v1.5:latest")

	// 4. 创建 Qdrant 向量存储
	store, err := vectorstore.NewQdrantStore(
		vectorstore.WithHost("106.55.106.158"),
		vectorstore.WithPort(6334),
		vectorstore.WithEmbedder(embedder),
	)
	if err != nil {
		t.Fatalf("Failed to create Qdrant store: %v", err)
	}

	// 5. 加载《孙子兵法》文档
	//file, err := os.Open("../三国演义.txt")
	//if err != nil {
	//	t.Fatalf("Failed to open szbf.txt: %v", err)
	//}
	//defer file.Close()
	//
	//loader := documentloader.New(file)
	//docs, err := loader.LoadSplit(splitter)
	//if err != nil {
	//	t.Fatalf("Failed to load szbf.txt: %v", err)
	//}

	// 6. 添加文档到向量存储
	//err = store.AddDocuments(context.Background(), "sgyy_collection", docs)
	//if err != nil {
	//	t.Fatalf("Failed to add documents to vector store: %v", err)
	//}

	// 7. 执行多个测试查询
	testQueries := []string{
		"三英战吕布是那几个人，结果如何？",
		"赤壁之战东吴诈降的是那位将军？",
		"三顾茅庐隆中对内容是什么",
		"解释一下白衣渡江",
	}

	for _, tc := range testQueries {
		t.Run(tc, func(t *testing.T) {
			results, err := store.SimilaritySearch(context.Background(), "sgyy_collection", tc, 5)
			if err != nil {
				t.Fatalf("Failed to perform similarity search for query '%s': %v", tc, err)
			}

			if len(results) == 0 {
				t.Fatal("No search results returned")
			}

			t.Logf("\nSearch results for query '%s':", tc)
			for i, result := range results {
				t.Logf("Result %d: Score=%.4f, Content: %s", i+1, result.Score, result.Metadata["content"])
			}
		})
	}
}
