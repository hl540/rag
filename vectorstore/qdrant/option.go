package qdrant

import (
	"crypto/tls"
	"github.com/hl540/rag/embedding"
	"google.golang.org/grpc"
)

type Option func(o *VectorStore)

func WithHost(host string) Option {
	return func(o *VectorStore) {
		o.config.Host = host
	}
}

func WithPort(port int) Option {
	return func(o *VectorStore) {
		o.config.Port = port
	}
}

func WithAPIKey(key string) Option {
	return func(o *VectorStore) {
		o.config.APIKey = key
	}
}

func WithUseTLS(use bool) Option {
	return func(o *VectorStore) {
		o.config.UseTLS = use
	}
}

func WithTLSConfig(config *tls.Config) Option {
	return func(o *VectorStore) {
		o.config.TLSConfig = config
	}
}

func WithGrpcOptions(opts []grpc.DialOption) Option {
	return func(o *VectorStore) {
		o.config.GrpcOptions = opts
	}
}

func WithSkipCompatibilityCheck(check bool) Option {
	return func(o *VectorStore) {
		o.config.SkipCompatibilityCheck = check
	}
}

func WithEmbedder(embedder embedding.Embedder) Option {
	return func(o *VectorStore) {
		o.embedder = embedder
	}
}
