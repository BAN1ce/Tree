package store

import (
	"github.com/BAN1ce/Tree/pkg"
	"github.com/BAN1ce/Tree/proto"
	"github.com/BAN1ce/Tree/state/api"
	proto2 "google.golang.org/protobuf/proto"
)

type Handle interface {
	HandleTopic
	HandleKeyStore
}

type HandleTopic interface {
	HandleSubRequest(req *proto.SubRequest) (response []byte, err error)
	HandleUnSubRequest(req *proto.UnSubRequest) (response []byte, err error)
}
type HandleKeyStore interface {
	HandlePutKeyRequest(req *proto.PutKeyRequest) (response []byte, err error)
	HandleDeleteKeyRequest(req *proto.DeleteKeyRequest) (response []byte, err error)
	HandlePutKeysRequest(req *proto.PutKeysRequest) (response []byte, err error)
}

func handleUpdateData(command []byte, handle Handle) ([]byte, error) {
	if len(command) < 1 {
		return nil, pkg.CommandLenError
	}
	var (
		body []byte
		err  error
	)
	switch command[0] {
	case api.SubCommand:
		var req proto.SubRequest
		if err = proto2.Unmarshal(command[1:], &req); err != nil {
			return nil, err
		}
		body, err = handle.HandleSubRequest(&req)
	case api.UnSubCommand:
		var req proto.UnSubRequest
		if err = proto2.Unmarshal(command[1:], &req); err != nil {
			return nil, err
		}
		body, err = handle.HandleUnSubRequest(&req)
	case api.PutKey:
		var req proto.PutKeyRequest
		if err = proto2.Unmarshal(command[1:], &req); err != nil {
			return nil, err
		}

		body, err = handle.HandlePutKeyRequest(&req)
	//case api.DefaultReadKey:
	//	var req proto.ReadKeyRequest
	//	if err = proto2.Unmarshal(command[1:], &req); err != nil {
	//		return nil, err
	//	}
	//	body, err = handle.HandleReadKeyRequest(&req)
	case api.DeleteKey:
		var req proto.DeleteKeyRequest
		if err = proto2.Unmarshal(command[1:], &req); err != nil {
			return nil, err
		}
		body, err = handle.HandleDeleteKeyRequest(&req)
	//case api.DefaultReadPrefixKey:
	//	var req proto.ReadPrefixKeyRequest
	//	if err = proto2.Unmarshal(command[1:], &req); err != nil {
	//		return nil, err
	//	}
	//	body, err = handle.HandleReadPrefixKeyRequest(&req)
	case api.PutKeys:
		var req proto.PutKeysRequest
		if err = proto2.Unmarshal(command[1:], &req); err != nil {
			return nil, err
		}
		body, err = handle.HandlePutKeysRequest(&req)
	//case api.ReadKeys:
	//	var req proto.ReadKeysRequest
	//	if err := proto2.Unmarshal(command[1:], &req); err != nil {
	//		return nil, err
	//	}
	//	body, err = handle.HandleReadKeysRequest(&req)
	default:
		return nil, pkg.ErrHandleTypeNotExist
	}
	if err != nil {
		return nil, err
	}
	//var rsp = make([]byte, 1, 1+len(body))
	//rsp[0] = command[0]
	//rsp = append(rsp, body...)
	return body, nil

}
