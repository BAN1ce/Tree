package topic

import (
	"fmt"
	"github.com/BAN1ce/Tree/pkg/util"
	"github.com/BAN1ce/Tree/proto"
	"github.com/BAN1ce/skyTree/pkg/utils"
	"github.com/eclipse/paho.golang/packets"
)

// ------------------------------------------------------- Share Topic -------------------------------------------------------//

type ShareTopic struct {
	root *proto.ShareTopicSubTree
}

func NewShareTopic() *ShareTopic {
	return &ShareTopic{
		root: newShareState(),
	}
}
func newShareState() *proto.ShareTopicSubTree {
	return &proto.ShareTopicSubTree{
		TreeRoot: &proto.ShareTreeNode{
			TopicSection: "/",
			Topic:        "/",
			TopicGroup:   make(map[string]*proto.ShareGroup),
			ChildNode:    nil,
		},
		Hash: make(map[string]*proto.ShareTopicTopicGroup),
	}
}

func (t *ShareTopic) handleShareTopic(clientID string, subOption *proto.SubOption) (err error) {
	var (
		groupName = subOption.GetShareGroup()
		topic     = subOption.GetTopic()
	)

	if util.HasWildcard(topic) {
		panic("implement me")
	}

	if groups, ok := t.root.Hash[topic]; ok {
		group, ok := groups.ShareGroups[groupName]

		if !ok {
			group = newShareGroups(topic, groupName)
		}

		if _, ok := group.Client[clientID]; !ok {
			group.Client[clientID] = subOption
			return nil
		} else {
			return fmt.Errorf("client %s already exists in group %s", clientID, groupName)
		}
	}

	t.root.Hash[topic] = newTopicGroup(topic)
	group := newShareGroups(topic, groupName)
	group.Client[clientID] = subOption
	t.root.Hash[topic].ShareGroups[groupName] = group
	return nil
}

func newTopicGroup(topic string) *proto.ShareTopicTopicGroup {
	return &proto.ShareTopicTopicGroup{
		ShareGroups: make(map[string]*proto.ShareGroup),
		Topic:       topic,
	}
}

func packetsToProtoSubOption(options *packets.SubOptions) *proto.SubOption {
	result := &proto.SubOption{
		QoS:               int32(options.QoS),
		NoLocal:           options.NoLocal,
		RetainAsPublished: options.RetainAsPublished,
		Topic:             options.Topic,
	}
	if utils.IsShareTopic(options.Topic) {
		result.Share = true
		result.ShareGroup, _ = utils.ParseShareTopic(options.Topic)
	}

	return result
}

func newShareGroups(topic, group string) *proto.ShareGroup {
	return &proto.ShareGroup{
		Client:    make(map[string]*proto.SubOption),
		Topic:     topic,
		GroupName: group,
	}
}
