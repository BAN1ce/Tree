package app

import (
	"context"
	"errors"
	"github.com/BAN1ce/Tree/inner/api"
	"github.com/BAN1ce/Tree/inner/api/request"
	"github.com/BAN1ce/Tree/inner/core"
	"github.com/BAN1ce/Tree/inner/server"
	"github.com/BAN1ce/Tree/logger"
	"github.com/BAN1ce/Tree/pkg"
	"github.com/BAN1ce/Tree/proto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	proto2 "google.golang.org/protobuf/proto"
	"sync"
	"time"
)

type App struct {
	mux          sync.RWMutex
	topicCluster *Cluster
	writeTimeout time.Duration
	gin          *gin.Engine
	State        *core.Core
}

func NewApp() *App {
	tmp := &App{
		State: core.NewState(),
	}
	tmp.Run()
	return tmp
}

func (app *App) Run() {
	if app.gin == nil {
		app.gin = gin.Default()
	}
	server.Route(app.gin, app.State)
}

func (app *App) StartTopicCluster(ctx context.Context, option ...Option) error {
	//option = append(option,cluster.WithStateMachine(store.NewState()))
	var (
		topicCluster = NewCluster(option...)
	)
	app.mux.Lock()
	defer app.mux.Unlock()
	app.topicCluster = topicCluster
	go func() {
		if err := topicCluster.Start(ctx); err != nil {
			logger.Logger.Fatal("store cluster start failed", zap.Error(err))
		}
	}()
	return nil
}

func (app *App) Subscribe(ctx context.Context, request *proto.SubRequest) (rsp *proto.SubResponse, err error) {
	data, err := api.EncodeRequestWriteCommand(request)
	if err != nil {
		return nil, err
	}
	if result, err := app.topicCluster.Write(ctx, data); err != nil {
		return nil, err
	} else {
		var rsp proto.SubResponse
		if err = proto2.Unmarshal(result.Data, &rsp); err != nil {
			return nil, err
		}
		return &rsp, nil
	}
}

func (app *App) UnSubscribe(ctx context.Context, request *proto.UnSubRequest) (rsp *proto.UnSubResponse, err error) {
	var (
		data []byte
	)
	data, err = api.EncodeRequestWriteCommand(request)
	if err != nil {
		return nil, err
	}
	if result, err := app.topicCluster.Write(ctx, data); err != nil {
		return nil, err
	} else {
		var rsp proto.UnSubResponse
		if err = proto2.Unmarshal(result.Data, &rsp); err != nil {
			return nil, err
		}
		return &rsp, nil
	}
}

func (app *App) PutKey(ctx context.Context, key, value string) error {
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
	if result, err := app.topicCluster.Write(ctx, data); err != nil {
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

func (app *App) ReadKey(ctx context.Context, key string) (value string, exists bool, err error) {
	var (
		req = &request.ReadKeyRequest{
			Key: key,
		}
		result interface{}
	)
	if result, err = app.topicCluster.Read(ctx, req); err != nil {
		return
	}
	if rsp, ok := result.(*request.ReadKeyResponse); ok {
		return rsp.Value, rsp.Exists, nil
	}
	err = pkg.ErrInvalidReadResponse
	return
}

func (app *App) DeleteKey(ctx context.Context, key string) error {
	var (
		req = &proto.DeleteKeyRequest{
			Key: key,
		}
	)
	data, err := api.EncodeRequestWriteCommand(req)
	if err != nil {
		return err
	}
	if result, err := app.topicCluster.Write(ctx, data); err != nil {
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

func (app *App) ReadPrefixKey(ctx context.Context, prefix string) (value map[string]string, err error) {
	var (
		req = &request.ReadPrefixKeyRequest{
			PrefixKey: prefix,
		}
		result interface{}
	)
	if result, err = app.topicCluster.Read(ctx, req); err != nil {
		return
	}
	if rsp, ok := result.(*request.ReadPrefixKeyResponse); ok {
		return rsp.Value, nil
	}
	return nil, pkg.ErrInvalidReadResponse
}

func (app *App) getWriteCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), app.writeTimeout)
}
