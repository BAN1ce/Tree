package topic

import (
	"fmt"
	"github.com/BAN1ce/Tree/proto"
	"github.com/BAN1ce/skyTree/pkg/utils"
)

// ------------------------------------------------------- Normal Topic -------------------------------------------------------//

type CommonTopic struct {
	root *proto.TopicSubTree
}

func NewCommonTopic() *CommonTopic {
	return &CommonTopic{
		root: newNormalState(),
	}
}

func newNormalState() *proto.TopicSubTree {
	return &proto.TopicSubTree{
		TreeRoot: &proto.TreeNode{
			TopicSection: "/",
			Topic:        "/",
			Clients:      make(map[string]*proto.SubOption, 100),
			ChildNode:    make(map[string]*proto.TreeNode, 1000000),
		},
		Hash: make(map[string]*proto.TopicClientsID),
	}

}

func (t *CommonTopic) createSub(clientID string, subOption *proto.SubOption) error {
	topicSection := utils.SplitTopic(subOption.Topic)
	if len(topicSection) == 0 {
		return fmt.Errorf("topic is empty")
	}
	treeRoot := t.root.TreeRoot
	for _, section := range topicSection {
		treeRoot.GetChildNode()
		if _, exists := treeRoot.ChildNode[section]; !exists {
			if treeRoot.ChildNode == nil {
				treeRoot.ChildNode = make(map[string]*proto.TreeNode)
			}
			// If the section does not exist, create a new TreeNode for it
			treeRoot.ChildNode[section] = &proto.TreeNode{
				TopicSection: section,
				Clients:      make(map[string]*proto.SubOption),
				ChildNode:    make(map[string]*proto.TreeNode),
			}
		}
		// Move to the child node representing the current section
		treeRoot = treeRoot.ChildNode[section]
	}
	treeRoot.Topic = subOption.Topic

	// Once all sections are processed, add the client's subscription option to the final node
	if _, exists := treeRoot.Clients[clientID]; !exists {
		if treeRoot.Clients == nil {
			treeRoot.Clients = make(map[string]*proto.SubOption, 100000)
		}
		treeRoot.Clients[clientID] = subOption
	}
	return nil
}

func (t *CommonTopic) deleteSub(topic string, clientID string) error {
	topicSection := utils.SplitTopic(topic)
	if len(topicSection) == 0 {
		return fmt.Errorf("topic is empty")
	}
	treeRoot := t.root.TreeRoot
	for _, section := range topicSection {
		if _, exists := treeRoot.ChildNode[section]; !exists {
			return nil
		}
		treeRoot = treeRoot.ChildNode[section]
	}
	delete(treeRoot.Clients, clientID)
	if len(treeRoot.Clients) == 0 {
		treeRoot.Clients = nil
	}
	return nil
}

func (t *CommonTopic) matchTopic(topic string) map[string]int32 {
	var (
		topicSections = utils.SplitTopic(topic)
		result        = make(map[string]int32)
		match         func(node *proto.TreeNode, sections []string, depth int)
	)

	match = func(node *proto.TreeNode, sections []string, depth int) {
		if depth == len(sections) {
			if node.Topic != "" {
				result[node.Topic] = t.nodeMaxQoS(node)
			}
			return
		}

		section := sections[depth]
		if node.ChildNode != nil {
			if child, ok := node.ChildNode[section]; ok {
				match(child, sections, depth+1)
			}
			if child, ok := node.ChildNode["+"]; ok && section != "/" {
				match(child, sections, depth+1)
			}
			if child, ok := node.ChildNode["#"]; ok && section != "/" {
				result[child.Topic] = t.nodeMaxQoS(child)
				return
			}
		}
	}

	if len(topicSections) > 0 {
		match(t.root.TreeRoot, topicSections, 0)
	}

	return result
}

func (t *CommonTopic) matchTopicForWildcard(wildcardTopic string) []string {
	var (
		result []string

		match func(node *proto.TreeNode, sections []string, depth int)
	)

	match = func(node *proto.TreeNode, sections []string, depth int) {
		if depth == len(sections) {
			if node.Topic != "" {
				result = append(result, node.Topic)
			}
			return
		}

		section := sections[depth]
		if section == "+" {
			for _, child := range node.ChildNode {
				if child.TopicSection != "#" && child.TopicSection != "+" {
					match(child, sections, depth+1)
				}
			}
		}

		if section == "#" {
			t.collectAllTopics(node, &result)
		}

		if section != "#" && section != "+" {
			if child, ok := node.ChildNode[section]; ok {
				match(child, sections, depth+1)
			}
		}

	}

	sections := utils.SplitTopic(wildcardTopic)
	match(t.root.TreeRoot, sections, 0)
	return result
}

func (t *CommonTopic) nodeMaxQoS(node *proto.TreeNode) int32 {
	var maxQoS int32
	for _, c := range node.Clients {
		if c.QoS > maxQoS {
			maxQoS = c.QoS
		}
	}
	return maxQoS
}

func (t *CommonTopic) collectAllTopics(node *proto.TreeNode, topics *[]string) {
	if node.Topic != "" {
		if node.TopicSection != "#" && node.TopicSection != "+" {
			*topics = append(*topics, node.Topic)
		}
	}
	for _, child := range node.ChildNode {
		t.collectAllTopics(child, topics)
	}
}
