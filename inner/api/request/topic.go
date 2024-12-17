package request

import (
	"github.com/BAN1ce/Tree/proto"
)

// --------------------------------------------------- sub topic ---------------------------------------------------//

type PostSubTopicRequest struct {
	Topic    []*proto.SubOption
	ClientID string
}

type PostSubTopicResponse struct {
	Result []error
}

// --------------------------------------------------- delete sub topic ---------------------------------------------------//

type DeleteSubTopicRequest struct {
	ClientID string
	Topic    []string
}

type DeleteSubTopicResponse struct {
	Result []error
}

// --------------------------------------------------- get all match topics ---------------------------------------------------//

type GetAllMatchTopicsRequest struct {
	Topic string
}

// GetAllMatchTopicsResponse is the response of GetAllMatchTopicsRequest
// eg: request topicName: "a/b/c", response topicName: "a/b/c" -> QoS, "a/b/# -> QoS, "a/#" -> QoS ......
// Topic: topicName -> QoS
type GetAllMatchTopicsResponse struct {
	RequestTopic string                 `json:"request_topic"`
	Topic        map[string]int32       `json:"topic"`
	ShareGroup   map[string]*ShareGroup `json:"share_group"`
}

type GetAllMatchTopicsForWildTopicRequest struct {
	Topic string
}

type GetAllMatchTopicsForWildTopicResponse struct {
	Topic []string `json:"topic"`
}

type ShareGroup struct {
	GroupName string           `json:"group_name"`
	Client    map[string]int32 `json:"client"`
}

// NewMatchSubTopicResponse is used to create a new GetAllMatchTopicsResponse
func NewMatchSubTopicResponse() *GetAllMatchTopicsResponse {
	return &GetAllMatchTopicsResponse{
		Topic:      make(map[string]int32),
		ShareGroup: make(map[string]*ShareGroup),
	}
}

// NewShareGroup is used to create a new ShareGroup
func NewShareGroup(groupName string) *ShareGroup {
	return &ShareGroup{
		GroupName: groupName,
		Client:    make(map[string]int32),
	}
}

// --------------------------------------------------- get client session ---------------------------------------------------//

type GetClientSessionRequest struct {
	ClientID string `json:"client_id"`
}

type GetClientSessionResponse struct {
	ClientID string `json:"client_id"`
}

// --------------------------------------------------- get topic info ---------------------------------------------------//

type GetTopicInfoRequest struct {
	Topic string `json:"topic"`
}

type GetTopicInfoResponse struct {
	Topic string `json:"topic"`
}
