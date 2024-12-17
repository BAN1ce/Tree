package topic

import (
	"context"
	"errors"
	"fmt"
	"github.com/BAN1ce/Tree/inner/api/request"
	"github.com/BAN1ce/skyTree/pkg/utils"
	"log/slog"
	"sync"
)

type SubTopic struct {
	normalSub  *CommonTopic
	shareTopic *ShareTopic
	mux        sync.RWMutex
}

func NewSubTopic() *SubTopic {
	return &SubTopic{
		normalSub:  NewCommonTopic(),
		shareTopic: NewShareTopic(),
	}
}

func (t *SubTopic) GetAllMatchTopics(ctx context.Context, req *request.GetAllMatchTopicsRequest) *request.GetAllMatchTopicsResponse {
	t.mux.RLock()
	defer t.mux.RUnlock()

	var (
		topic    = req.Topic
		response = request.NewMatchSubTopicResponse()
	)
	response.RequestTopic = req.Topic

	matchTopic := t.normalSub.matchTopic(topic)
	for topic, maxQoS := range matchTopic {
		response.Topic[topic] = maxQoS
	}

	if groups, ok := t.shareTopic.root.Hash[topic]; ok {
		for _, group := range groups.ShareGroups {
			var (
				shareGroup = request.NewShareGroup(group.GroupName)
			)

			response.Topic[fmt.Sprintf("$share/%s/%s", shareGroup.GroupName, topic)] = 1

			for clientID, subOption := range group.Client {
				shareGroup.Client[clientID] = subOption.GetQoS()
			}
			response.ShareGroup[group.GroupName] = shareGroup
		}
	}

	return response

}

func (t *SubTopic) PostSubTopic(ctx context.Context, req *request.PostSubTopicRequest) *request.PostSubTopicResponse {
	t.mux.Lock()
	defer t.mux.Unlock()

	var (
		response = &request.PostSubTopicResponse{}
		errs     error
		err      error
	)

	for _, subOption := range req.Topic {
		var (
			topic    = subOption.Topic
			clientID = req.ClientID
		)
		if utils.IsShareTopic(topic) {
			err = t.shareTopic.handleShareTopic(clientID, subOption)
			if err != nil {
				slog.Error("handle share sub request error", err.Error(), req)
				errs = errors.Join(errs, err)
			}
			response.Result = append(response.Result, err)
			continue
		}

		err = t.normalSub.createSub(req.ClientID, subOption)
		if err != nil {
			slog.Error("create wildcard sub error", err.Error(), req)
			errs = errors.Join(errs, err)
		}
		response.Result = append(response.Result, err)

	}
	return response

}

func (t *SubTopic) DeleteSubTopic(ctx context.Context, req *request.DeleteSubTopicRequest) *request.DeleteSubTopicResponse {
	var (
		err      error
		response = &request.DeleteSubTopicResponse{}
	)
	t.mux.Lock()
	defer t.mux.Unlock()

	for _, topic := range req.Topic {
		if utils.IsShareTopic(topic) {
			// TODO: delete share topic
			panic("implement me")
		}
		err = t.normalSub.deleteSub(topic, req.ClientID)
		if err != nil {
			slog.Error("delete wildcard sub error", err.Error(), req)
		}
		response.Result = append(response.Result, err)
	}
	return response
}

func (t *SubTopic) GetMatchTopicForWildcardTopic(ctx context.Context, req *request.GetAllMatchTopicsForWildTopicRequest) *request.GetAllMatchTopicsForWildTopicResponse {
	t.mux.RLock()
	defer t.mux.RUnlock()
	var (
		response = &request.GetAllMatchTopicsForWildTopicResponse{}
	)
	response.Topic = t.normalSub.matchTopicForWildcard(req.Topic)
	return response
}
