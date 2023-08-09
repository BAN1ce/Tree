package api

import (
	"github.com/BAN1ce/Tree/pkg"
	"github.com/BAN1ce/Tree/proto"
	proto2 "google.golang.org/protobuf/proto"
)

const (
	SubCommand = byte(iota)
	UnSubCommand
	PutKey
	ReadKey
	ReadPrefixKey
	DeleteKey
	PutKeys
	ReadKeys
)

func EncodeRequestWriteCommand(c interface{}) (command []byte, err error) {
	if req, ok := c.(proto2.Message); ok {
		command, err = proto2.Marshal(req)
		if err != nil {
			return
		}
	} else {
		err = pkg.ErrCommandTypeNotExist
		return
	}
	switch c.(type) {
	case *proto.SubRequest:
		command = append([]byte{SubCommand}, command...)
		return
	case *proto.UnSubRequest:
		command = append([]byte{UnSubCommand}, command...)
		return

	case *proto.PutKeyRequest:
		command = append([]byte{PutKey}, command...)
		return
	case *proto.ReadKeyRequest:
		command = append([]byte{ReadKey}, command...)
	case *proto.DeleteKeyRequest:
		command = append([]byte{DeleteKey}, command...)
	case *proto.ReadPrefixKeyRequest:
		command = append([]byte{ReadPrefixKey}, command...)
	default:
		return nil, pkg.ErrHandleTypeNotExist
	}
	return
}

type MatchTopicRequest struct {
	Topic string
}

type ClientInfo struct {
	ClientID    string
	QoS         int32
	NodeAddress string
}
type MatchTopicResponse struct {
	Client []*ClientInfo
}
