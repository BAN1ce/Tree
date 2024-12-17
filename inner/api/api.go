package api

import (
	"context"
	"github.com/BAN1ce/Tree/inner/api/request"
)

type API interface {
	GetAllMatchTopics(ctx context.Context, req *request.GetAllMatchTopicsRequest) *request.GetAllMatchTopicsResponse
	GetMatchTopicForWildcardTopic(ctx context.Context, req *request.GetAllMatchTopicsForWildTopicRequest) *request.GetAllMatchTopicsForWildTopicResponse
}

type SessionAPI interface {
	PutKey(ctx context.Context, key, value string) error

	ReadKey(ctx context.Context, key string) (string, bool, error)

	ReadPrefixKey(ctx context.Context, prefix string) (map[string]string, error)

	DeleteKey(ctx context.Context, key string) error

	DeletePrefixKey(ctx context.Context, prefix string) (map[string]string, error)
}

type TopicAPI interface {
	GetAllMatchTopics(ctx context.Context, req *request.GetAllMatchTopicsRequest) *request.GetAllMatchTopicsResponse
	PostSubTopic(ctx context.Context, req *request.PostSubTopicRequest) *request.PostSubTopicResponse
	DeleteSubTopic(ctx context.Context, req *request.DeleteSubTopicRequest) *request.DeleteSubTopicResponse
	GetMatchTopicForWildcardTopic(ctx context.Context, req *request.GetAllMatchTopicsForWildTopicRequest) *request.GetAllMatchTopicsForWildTopicResponse
}

type StatusAPI interface {
}
