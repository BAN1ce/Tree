package core

import (
	"context"
	"github.com/BAN1ce/Tree/inner/api/request"
	session2 "github.com/BAN1ce/Tree/inner/core/session"
	"github.com/BAN1ce/Tree/inner/core/topic"
	"github.com/BAN1ce/Tree/proto"
)

type Core struct {
	subTopic *topic.SubTopic
	session  *session2.Session
}

func NewState() *Core {
	return &Core{
		session:  session2.NewSession(),
		subTopic: topic.NewSubTopic(),
	}
}

func (c *Core) HandleSubRequest(req *proto.SubRequest) *request.PostSubTopicResponse {
	var (
		apiReq = request.PostSubTopicRequest{}
	)

	apiReq.ClientID = req.GetClientID()
	for _, option := range req.GetTopics() {
		apiReq.Topic = append(apiReq.Topic, option)
	}

	return c.subTopic.PostSubTopic(context.TODO(), &apiReq)
}

func (c *Core) HandleUnSubRequest(req *proto.UnSubRequest) *request.DeleteSubTopicResponse {
	var (
		apiReq = &request.DeleteSubTopicRequest{
			ClientID: req.GetClientID(),
			Topic:    req.GetTopics(),
		}
	)
	return c.subTopic.DeleteSubTopic(context.TODO(), apiReq)
}

func (c *Core) GetAllMatchTopics(ctx context.Context, req *request.GetAllMatchTopicsRequest) *request.GetAllMatchTopicsResponse {
	return c.subTopic.GetAllMatchTopics(ctx, req)
}

func (c *Core) GetMatchTopicForWildcardTopic(ctx context.Context, req *request.GetAllMatchTopicsForWildTopicRequest) *request.GetAllMatchTopicsForWildTopicResponse {
	return c.subTopic.GetMatchTopicForWildcardTopic(ctx, req)

}
