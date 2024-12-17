package session

import (
	"context"
	"errors"
	"fmt"
	"github.com/BAN1ce/Tree/proto"
	"github.com/BAN1ce/skyTree/logger"
	"strings"
	"sync"
)

type Session struct {
	*Node
	mux sync.RWMutex
}

func NewSession() *Session {
	return &Session{
		Node: &Node{
			Node: NewNode(),
		},
	}
}

func (s *Session) PutKey(ctx context.Context, key, value string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.Node.PutKey(key, value)
}

func (s *Session) ReadKey(ctx context.Context, key string) (string, bool, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.Node.ReadKey(key)
}

func (s *Session) ReadPrefixKey(ctx context.Context, prefix string) (map[string]string, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.Node.PrefixSearch(prefix)
}

func (s *Session) DeleteKey(ctx context.Context, key string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.Node.DeleteKey(key)
}

func (s *Session) DeletePrefixKey(ctx context.Context, prefix string) (map[string]string, error) {
	s.mux.Lock()
	s.mux.Unlock()

	var (
		errs error
	)

	result, err := s.Node.PrefixSearch(prefix)
	if err != nil {
		return nil, err
	}
	for k := range result {
		if err = s.Node.DeleteKey(k); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return result, errs
}

//---------------------------------Prefix Tree Node---------------------------------//

// Node represents a node in the prefix tree.go.
// It contains a key-value pair and a map of child nodes.
// The key-value pair is stored in the leaf node.
// The map of child nodes is stored in the non-leaf node.
type Node struct {
	*proto.Node
}

func NewNode() *proto.Node {
	return &proto.Node{
		ChildNode: make(map[string]*proto.Node),
	}
}

// PutKey inserts a key-value pair into the prefix tree.go.
func (t *Node) PutKey(key, value string) error {
	if key == "" {
		return errors.New("empty key is not allowed")
	}

	keys := strings.Split(key, "/")
	node := t.Node
	for _, k := range keys {
		child, found := node.ChildNode[k]
		if !found {
			child = NewNode()
			child.Key = k
			if node.FullKey == "" {
				child.FullKey = k
			} else {
				child.FullKey = node.FullKey + "/" + k
			}
			node.ChildNode[k] = child
		}
		node = child
	}

	node.Value = value
	return nil
}

// ReadKey retrieves the value associated with the given key from the prefix tree.go.
func (t *Node) ReadKey(key string) (string, bool, error) {
	if key == "" {
		return "", false, errors.New("empty key is not allowed")
	}

	keys := strings.Split(key, "/")
	node := t.Node
	for _, k := range keys {
		child, found := node.ChildNode[k]
		if !found {
			return "", false, nil
		}
		node = child
	}
	value := node.Value
	return value, true, nil
}

// DeleteKey deletes a key-value pair from the prefix tree.go.
func (t *Node) DeleteKey(key string) error {
	if key == "" {
		return errors.New("empty key is not allowed")
	}

	keys := strings.Split(key, "/")
	parent, node := t.Node, t.Node
	for _, k := range keys {
		child, found := node.ChildNode[k]
		if !found {
			logger.Logger.Info().Str("key", key).Msg("sub tree delete key not found")
			return nil
		}
		parent, node = node, child
	}

	// Delete the node and release the map if it has no children.
	delete(parent.ChildNode, node.Key)
	if len(parent.ChildNode) == 0 {
		parent.ChildNode = map[string]*proto.Node{}
	}

	return nil
}

// PrefixSearch searches for all key-value pairs with keys that have the given prefix.
func (t *Node) PrefixSearch(prefix string) (map[string]string, error) {
	if prefix == "" {
		return nil, fmt.Errorf("empty prefix is not allowed")
	}

	result := make(map[string]string)
	keys := strings.Split(prefix, "/")
	node := t.Node
	for _, k := range keys {
		child, found := node.ChildNode[k]
		if !found {
			return nil, nil
		}
		node = child
	}

	t.traverse(node, prefix, result)

	return result, nil
}

// traverse recursively traverses the subtree from the given node and collects key-value pairs with the given prefix.
func (t *Node) traverse(node *proto.Node, prefix string, result map[string]string) {
	if node.Value != "" {
		result[node.FullKey] = node.Value
	}
	for _, child := range node.ChildNode {
		t.traverse(child, prefix, result)
	}
}
