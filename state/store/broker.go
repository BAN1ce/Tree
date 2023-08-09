package store

import (
	"errors"
	"github.com/BAN1ce/Tree/pkg"
	"github.com/BAN1ce/Tree/pkg/util"
	"github.com/BAN1ce/Tree/proto"
	"github.com/BAN1ce/Tree/state/api"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/lni/dragonboat/v3/statemachine"
	"go.uber.org/zap"
	proto2 "google.golang.org/protobuf/proto"
	"io"
	"strings"
	"sync"
)

type State struct {
	state   *proto.TopicSubTree
	session *Node
	mux     sync.RWMutex
}

func NewState() *State {
	return &State{
		state: &proto.TopicSubTree{
			TreeRoot: &proto.TreeNode{
				TopicSection: "/",
				Topic:        "/",
				Clients:      make(map[string]int32),
				ChildNode:    nil,
			},
			Hash: make(map[string]*proto.TopicClientsID),
		},
		session: &Node{
			Node: NewNode(),
		},
	}
}
func (t *State) Update(bytes []byte) (statemachine.Result, error) {
	t.mux.Lock()
	defer t.mux.Unlock()
	var (
		result statemachine.Result
	)
	data, err := handleUpdateData(bytes, t)
	if err != nil {
		return result, err
	} else {
		result.Data = data
		return result, nil
	}
}

func (t *State) Lookup(i interface{}) (interface{}, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	switch req := i.(type) {
	case *api.MatchTopicRequest:
		return t.MatchTopic(req), nil
	case *api.ReadKeyRequest:
		v, exist, err := t.ReadKey(req.Key)
		if err != nil {
			return nil, err
		}
		return &api.ReadKeyResponse{
			Key:    req.Key,
			Value:  v,
			Exists: exist,
		}, nil
	case *api.ReadPrefixKeyRequest:
		v, err := t.ReadWithPrefix(req.PrefixKey)
		if err != nil {
			return nil, err
		}
		return &api.ReadPrefixKeyResponse{
			PrefixKey: req.PrefixKey,
			Value:     v,
		}, nil

	default:
		return nil, pkg.ErrHandleTypeNotExist
	}
}

func (t *State) SaveSnapshot(writer io.Writer, collection statemachine.ISnapshotFileCollection, i <-chan struct{}) error {
	t.mux.Lock()
	defer t.mux.Unlock()
	body, err := proto2.Marshal(t.state)
	if err != nil {
		return err
	}
	_, err = writer.Write(body)
	return err
}

func (t *State) RecoverFromSnapshot(reader io.Reader, files []statemachine.SnapshotFile, i <-chan struct{}) error {
	t.mux.Lock()
	defer t.mux.Unlock()
	var (
		data, err = io.ReadAll(reader)
	)
	if err != nil {
		return err
	}
	return proto2.Unmarshal(data, t.state)
}

func (t *State) Close() error {
	t.mux.Lock()
	defer t.mux.Unlock()
	return nil
}

func (t *State) HandleSubRequest(req *proto.SubRequest) (response []byte, err error) {
	for topic, info := range req.GetTopics() {
		if util.HasWildcard(topic) {
			panic("implement me")
		}
		if clients, ok := t.state.Hash[topic]; ok {
			clients.Clients[req.ClientID] = info.QoS
		} else {
			t.state.Hash[topic] = &proto.TopicClientsID{
				Clients: map[string]int32{
					req.GetClientID(): info.QoS,
				},
			}
		}
	}
	// TODO: implement me response
	return nil, nil
}

func (t *State) HandleUnSubRequest(req *proto.UnSubRequest) (response []byte, err error) {
	for _, topic := range req.GetTopics() {
		if util.HasWildcard(topic) {
			panic("implement me")
		}
		if clients, ok := t.state.Hash[topic]; ok {
			delete(clients.Clients, req.GetClientID())
			if len(clients.Clients) == 0 {
				delete(t.state.Hash, topic)
			}
		}
	}
	// TODO: implement me response
	return nil, nil
}

func (t *State) MatchTopic(req *api.MatchTopicRequest) *api.MatchTopicResponse {
	var (
		response api.MatchTopicResponse
	)
	if util.HasWildcard(req.Topic) {
		panic("implement me")
	}
	if clients, ok := t.state.Hash[req.Topic]; ok {
		for clientID, qos := range clients.Clients {
			response.Client = append(response.Client, &api.ClientInfo{
				ClientID: clientID,
				QoS:      qos,
			})
		}
	}
	return &response
}

func (t *State) HandlePutKeyRequest(req *proto.PutKeyRequest) (response []byte, err error) {
	var (
		resp proto.PutKeyResponse
	)
	if err = t.session.PutKey(req.GetKey(), req.GetValue()); err != nil {
		resp.Message = err.Error()
		return
	}
	resp.Success = true
	response, err = proto2.Marshal(&resp)
	return
}

func (t *State) HandleDeleteKeyRequest(req *proto.DeleteKeyRequest) (response []byte, err error) {
	var (
		resp proto.DeleteKeyResponse
	)
	if err = t.session.DeleteKey(req.GetKey()); err == nil {
		resp.Success = true
	} else {
		resp.Message = err.Error()
	}
	response, err = proto2.Marshal(&resp)
	return
}

func (t *State) HandlePutKeysRequest(req *proto.PutKeysRequest) (response []byte, err error) {
	var (
		resp   = proto.PutKeysResponse{}
		tmpErr error
	)
	for k, v := range req.Value {
		tmpErr = t.session.PutKey(k, v)
		if tmpErr != nil {
			err = errors.Join(err, tmpErr)
		}
	}
	if err == nil {
		resp.Success = true
	}

	response, tmpErr = proto2.Marshal(&resp)
	err = errors.Join(err, tmpErr)
	return
}

func (t *State) ReadKey(key string) (string, bool, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.session.ReadKey(key)
}

func (t *State) ReadWithPrefix(prefix string) (map[string]string, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.session.PrefixSearch(prefix)
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
			child.FullKey = node.FullKey + "/" + k
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

	return node.Value, true, nil
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
			logger.Logger.Info("key not found", zap.String("key", key))
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
		return nil, errors.New("empty prefix is not allowed")
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

func (t *Node) ReadWildcardKey(wildcardKey string) (map[string]string, error) {
	if wildcardKey == "" {
		return nil, errors.New("empty wildcard key is not allowed")
	}

	result := make(map[string]string)
	t.readWildcard(t.Node, wildcardKey, "", result)

	return result, nil
}

// readWildcard recursively searches for key-value pairs that match the given wildcard key.
func (t *Node) readWildcard(node *proto.Node, wildcardKey, currentKey string, result map[string]string) {
	if wildcardKey == "" {
		if node.Value != "" {
			result[node.FullKey] = node.Value
		}
		return
	}

	// Check for '+'
	if wildcardKey[0] == '+' {
		for _, child := range node.ChildNode {
			t.readWildcard(child, wildcardKey[1:], currentKey+child.Key+"/", result)
		}
		return
	}

	// Check for '*'
	if wildcardKey[0] == '*' {
		for _, child := range node.ChildNode {
			t.readWildcard(child, wildcardKey, currentKey+child.Key+"/", result)
		}
		t.readWildcard(node, wildcardKey[1:], currentKey, result)
		return
	}

	// Normal character
	if child, found := node.ChildNode[string(wildcardKey[0])]; found {
		t.readWildcard(child, wildcardKey[1:], currentKey+child.Key+"/", result)
	}
}
