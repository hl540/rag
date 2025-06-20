package vectorstore

import (
	"crypto/tls"
	"github.com/hl540/rag/embedding"
	"google.golang.org/grpc"
)

type QdrantOption func(o *QdrantStore)

func WithHost(host string) QdrantOption {
	return func(o *QdrantStore) {
		o.config.Host = host
	}
}

func WithPort(port int) QdrantOption {
	return func(o *QdrantStore) {
		o.config.Port = port
	}
}

func WithAPIKey(key string) QdrantOption {
	return func(o *QdrantStore) {
		o.config.APIKey = key
	}
}

func WithUseTLS(use bool) QdrantOption {
	return func(o *QdrantStore) {
		o.config.UseTLS = use
	}
}

func WithTLSConfig(config *tls.Config) QdrantOption {
	return func(o *QdrantStore) {
		o.config.TLSConfig = config
	}
}

func WithGrpcOptions(opts []grpc.DialOption) QdrantOption {
	return func(o *QdrantStore) {
		o.config.GrpcOptions = opts
	}
}

func WithSkipCompatibilityCheck(check bool) QdrantOption {
	return func(o *QdrantStore) {
		o.config.SkipCompatibilityCheck = check
	}
}

func WithEmbedder(embedder embedding.Embedder) QdrantOption {
	return func(o *QdrantStore) {
		o.embedder = embedder
	}
}
