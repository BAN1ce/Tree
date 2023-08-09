package app

import (
	"context"
	"errors"
	"github.com/BAN1ce/Tree/logger"
	"github.com/BAN1ce/Tree/pkg"
	"github.com/BAN1ce/Tree/proto"
	"github.com/BAN1ce/Tree/state/api"
	"go.uber.org/zap"
	proto2 "google.golang.org/protobuf/proto"
	"sync"
	"time"
)

type App struct {
	mux          sync.RWMutex
	topicCluster *Cluster
	writeTimeout time.Duration
}

func NewApp() *App {
	return &App{}
}

func (a *App) StartTopicCluster(ctx context.Context, option ...Option) error {
	//option = append(option,cluster.WithStateMachine(store.NewState()))
	var (
		topicCluster = NewCluster(option...)
	)
	a.mux.Lock()
	defer a.mux.Unlock()
	a.topicCluster = topicCluster
	go func() {
		if err := topicCluster.Start(ctx); err != nil {
			logger.Logger.Fatal("store cluster start failed", zap.Error(err))
		}
	}()
	return nil
}

func (a *App) Subscribe(ctx context.Context, request *proto.SubRequest) (rsp *proto.SubResponse, err error) {
	data, err := api.EncodeRequestWriteCommand(request)
	if err != nil {
		return nil, err
	}
	if result, err := a.topicCluster.Write(ctx, data); err != nil {
		return nil, err
	} else {
		var rsp proto.SubResponse
		if err = proto2.Unmarshal(result.Data, &rsp); err != nil {
			return nil, err
		}
		return &rsp, nil
	}
}

func (a *App) UnSubscribe(ctx context.Context, request *proto.UnSubRequest) (rsp *proto.UnSubResponse, err error) {
	var (
		data []byte
	)
	data, err = api.EncodeRequestWriteCommand(request)
	if err != nil {
		return nil, err
	}
	if result, err := a.topicCluster.Write(ctx, data); err != nil {
		return nil, err
	} else {
		var rsp proto.UnSubResponse
		if err = proto2.Unmarshal(result.Data, &rsp); err != nil {
			return nil, err
		}
		return &rsp, nil
	}
}

func (a *App) PutKey(ctx context.Context, key, value string) error {
	var (
		req = &proto.PutKeyRequest{
			Key:   key,
			Value: value,
		}
	)
	data, err := api.EncodeRequestWriteCommand(req)
	if err != nil {
		return err
	}
	// FIXME ctx
	if result, err := a.topicCluster.Write(ctx, data); err != nil {
		return err
	} else {
		var rsp proto.PutKeyResponse
		if err = proto2.Unmarshal(result.Data, &rsp); err != nil {
			return err
		}
		if rsp.Success {
			return nil
		}
		return errors.New(rsp.Message)
	}

}

func (a *App) ReadKey(ctx context.Context, key string) (value string, exists bool, err error) {
	var (
		req = &api.ReadKeyRequest{
			Key: key,
		}
		result interface{}
	)
	if result, err = a.topicCluster.Read(ctx, req); err != nil {
		return
	}
	if rsp, ok := result.(*api.ReadKeyResponse); ok {
		return rsp.Value, rsp.Exists, nil
	}
	err = pkg.ErrInvalidReadResponse
	return
}

func (a *App) DeleteKey(ctx context.Context, key string) error {
	var (
		req = &proto.DeleteKeyRequest{
			Key: key,
		}
	)
	data, err := api.EncodeRequestWriteCommand(req)
	if err != nil {
		return err
	}
	if result, err := a.topicCluster.Write(ctx, data); err != nil {
		return err
	} else {
		var rsp proto.DeleteKeyResponse
		if err = proto2.Unmarshal(result.Data, &rsp); err != nil {
			return err
		}
		if rsp.Success {
			return nil
		}
		return errors.New(rsp.Message)
	}
}

func (a *App) ReadPrefixKey(ctx context.Context, prefix string) (value map[string]string, err error) {
	var (
		req = &api.ReadPrefixKeyRequest{
			PrefixKey: prefix,
		}
		result interface{}
	)
	if result, err = a.topicCluster.Read(ctx, req); err != nil {
		return
	}
	if rsp, ok := result.(*api.ReadPrefixKeyResponse); ok {
		return rsp.Value, nil
	}
	return nil, pkg.ErrInvalidReadResponse
}

func (a *App) MatchTopic(ctx context.Context, req *api.MatchTopicRequest) (*api.MatchTopicResponse, error) {
	result, err := a.topicCluster.Read(ctx, req)
	if err != nil {
		return nil, err
	}
	if rsp, ok := result.(*api.MatchTopicResponse); ok {
		return rsp, nil
	}
	return nil, pkg.ErrInvalidReadResponse
}

func (a *App) getWriteCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), a.writeTimeout)
}
