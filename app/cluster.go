package app

import (
	"context"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
	"github.com/lni/dragonboat/v3/statemachine"
	"go.uber.org/zap"
	"strconv"
)

type Option func(*Cluster)

func WithInitMember(mem map[uint64]dragonboat.Target) Option {
	return func(cluster *Cluster) {
		cluster.initMember = mem
	}
}

func WithJoin(join bool) Option {
	return func(cluster *Cluster) {
		cluster.join = join
	}
}

func WithConfig(cfg config.Config) Option {
	return func(cluster *Cluster) {
		cluster.cfg = cfg
	}
}

func WithStateMachine(state statemachine.IStateMachine) Option {
	return func(cluster *Cluster) {
		cluster.state = state
	}
}
func WithNodeConfig(nodeConfig config.NodeHostConfig) Option {
	return func(cluster *Cluster) {
		cluster.nodeConfig = nodeConfig
	}
}

type Cluster struct {
	initMember map[uint64]dragonboat.Target
	join       bool
	ctx        context.Context
	cancel     context.CancelFunc
	cfg        config.Config
	nodeConfig config.NodeHostConfig
	state      statemachine.IStateMachine
	clusterID  uint64
	node       *dragonboat.NodeHost
}

func NewCluster(options ...Option) *Cluster {
	cluster := &Cluster{
		initMember: make(map[uint64]dragonboat.Target),
		join:       false,
		cfg:        config.Config{},
		nodeConfig: config.NodeHostConfig{},
		state:      nil,
		clusterID:  0,
		node:       nil,
	}
	for _, option := range options {
		option(cluster)
	}
	cluster.clusterID = cluster.cfg.ClusterID
	return cluster

}

func (c *Cluster) Start(ctx context.Context) error {
	var err error
	c.ctx, c.cancel = context.WithCancel(ctx)
	c.node, err = dragonboat.NewNodeHost(c.nodeConfig)
	if err != nil {
		panic(err)
	}
	if err := c.node.StartCluster(c.initMember, c.join, func(u uint64, u2 uint64) statemachine.IStateMachine {
		return c.state
	}, c.cfg); err != nil {
		logger.Logger.Error("failed to start cluster", zap.Error(err))
		return err
	}
	<-c.ctx.Done()
	return c.close()
}

func (c *Cluster) Close() error {
	c.cancel()
	return nil
}

func (c *Cluster) close() error {
	nodeID, err := strconv.ParseUint(c.node.ID(), 10, 64)
	if err != nil {
		logger.Logger.Error("failed to parse node id", zap.Error(err))
		return err
	}
	if err := c.node.StopNode(c.clusterID, nodeID); err != nil {
		logger.Logger.Error("failed to stop node", zap.Error(err))
		return err
	}
	return err
}

func (c *Cluster) Write(ctx context.Context, data []byte) (statemachine.Result, error) {
	var session = c.node.GetNoOPSession(c.clusterID)
	return c.node.SyncPropose(ctx, session, data)
}

func (c *Cluster) Read(ctx context.Context, query interface{}) (interface{}, error) {
	return c.node.SyncRead(ctx, c.clusterID, query)
}
