package test

import (
	"context"
	"github.com/google/uuid"
	ollama2 "github.com/hl540/rag/embedding/ollama"
	"github.com/hl540/rag/llm/ollama"
	"github.com/hl540/rag/vectorstore"
	"github.com/hl540/rag/vectorstore/qdrant"
	"github.com/joho/godotenv"
	"github.com/ollama/ollama/api"
	"log"
	"testing"
)

func TestEmbedText(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("加载 .env 失败: %v", err)
	}

	llm, err := ollama.New("http://127.0.0.1:11434", "snowflake-arctic-embed:33m")
	if err != nil {
		log.Fatalf("llm 客户端创建失败: %v", err)
	}
	ctx := context.Background()
	embed, err := llm.CreateEmbedding(ctx, "hello")
	if err != nil {
		log.Fatalf("llm embedding 失败: %v", err)
	}
	t.Log(embed)
}

func TestVectorStoreAddDocument(t *testing.T) {
	ollama, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatalf("ollama 连接失败: %v", err)
	}

	vectorStore, err := qdrant.New(
		qdrant.WithHost("106.55.106.158"),
		qdrant.WithPort(6334),
		qdrant.WithEmbedder(ollama2.New(ollama, "snowflake-arctic-embed:33m")),
	)
	if err != nil {
		log.Fatalf("qdrant 连接失败: %v", err)
	}
	ctx := context.Background()

	texts := []*vectorstore.Document{
		{
			Id:       uuid.NewString(),
			Text:     "特朗普24小时内连出三记重拳：威胁苹果海外生产加征25%关税，对欧盟商品祭出50%惩罚性关税，同时签署核能复兴行政令。这场以\"美国优先\"为名的政策风暴，或将推高全球通胀、引爆贸易战，并埋下能源安全隐忧。",
			Metadata: make(map[string]any),
		},
		{
			Id:       uuid.NewString(),
			Text:     "在讲话中，钟南山提到了在今年2月份因流感离世的女明星大S徐熙媛，并表示很遗憾。他强调，近期流感发病率已经在下降，要高度重视流感的及时治疗，特别对于老年人来说，要保证在48小时内及时接受治疗。“药物的选择也很重要，国内近期有两款新药上市，目前中国在研发流感药物及疫苗方面并不比国外差。",
			Metadata: make(map[string]any),
		},
		{
			Id:       uuid.NewString(),
			Text:     "相关报道称，近日，四川绵阳市商务局印发《绵阳市提振消费专项行动2025年工作清单》，在备受关注的落实休假政策中提到，落实年休假应休尽休和带薪休假政策，鼓励企业弹性调休，推广夫妻共享。还试行4.5天弹性工作制，鼓励有条件的地区推行“周五下午与周末结合”的2.5天休假模式。",
			Metadata: make(map[string]any),
		},
		{
			Id:       uuid.NewString(),
			Text:     "5月23日晚，格力董事长董明珠和前秘书孟羽童一同现身“格力明珠精选”直播间进行直播带货。这是孟羽童离开格力电器后两人首次同框“破冰”。",
			Metadata: make(map[string]any),
		},
	}
	for i := range texts {
		texts[i].Metadata[vectorstore.ContentKey] = texts[i].Text
	}
	err = vectorStore.AddDocuments(ctx, "test", texts)
	if err != nil {
		log.Fatalf("vector 新增失败: %v", err)
	}
	t.Log("success")
}

func TestVectorStoreSimilaritySearch(t *testing.T) {
	ollama, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatalf("ollama 连接失败: %v", err)
	}

	vectorStore, err := qdrant.New(
		qdrant.WithHost("106.55.106.158"),
		qdrant.WithPort(6334),
		qdrant.WithEmbedder(ollama2.New(ollama, "snowflake-arctic-embed:33m")),
	)
	if err != nil {
		log.Fatalf("qdrant 连接失败: %v", err)
	}
	ctx := context.Background()
	search, err := vectorStore.SimilaritySearch(ctx, "test", "格力董事长董明珠和前秘书孟羽童一同现身", 5)
	if err != nil {
		log.Fatalf("vector 查询失败: %v", err)
	}
	for _, doc := range search {
		t.Logf("ID:%s, Score:%f, Metadata:%+v", doc.ID, doc.Score, doc.Metadata)
	}
}
