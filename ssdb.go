package main

import (
	"errors"
	"log"
	"strings"

	"github.com/seefan/gossdb"
)

type SSDBPool struct {
	logger     *log.Logger
	Config     *Config
	MasterPool *gossdb.Connectors
	PoolMap    map[string]*gossdb.Connectors
}

func NewSSDBPool(logger *log.Logger, config *Config) *SSDBPool {
	return &SSDBPool{
		logger:     logger,
		Config:     config,
		MasterPool: nil,
		PoolMap:    make(map[string]*gossdb.Connectors),
	}
}

func (p *SSDBPool) Start() {
	defaultArgs := p.Config.Args
	conn, err := newPool(p.Config.Master, defaultArgs)
	if err != nil {
		p.logger.Println(err)
	}
	p.MasterPool = conn
	for _, node := range p.Config.Nodes {
		c, err := newPool(node, defaultArgs)
		if err != nil {
			p.logger.Println(err)
			continue
		}
		p.PoolMap[node.Route] = c
	}
	p.logger.Println("SSDBPool start success")
	p.logger.Println(p.MasterPool, p.PoolMap)
}

func newPool(node Node, defaultArgs Args) (*gossdb.Connectors, error) {
	if !node.Open {
		return nil, errors.New("node not open")
	}
	args := defaultArgs
	if node.OwnArgs {
		args = node.Args
	}
	return gossdb.NewPool(&gossdb.Config{
		Host:             node.Host,
		Port:             node.Port,
		GetClientTimeout: args.GetClientTimeout,
		MaxPoolSize:      args.MaxPoolSize,
		MinPoolSize:      args.MinPoolSize,
		AcquireIncrement: args.AcquireIncrement,
		MaxIdleTime:      args.MaxIdleTime,
		MaxWaitSize:      args.MaxWaitSize,
		HealthSecond:     args.HealthSecond,
	}, node.Password)
}

func (p *SSDBPool) Do(args []string) ([]string, error) {
	pool := p.routePool(args)
	if pool == nil {
		return nil, errors.New("route failed")
	}
	c, err := pool.NewClient()
	if err != nil {
		p.logger.Println("ssdb new client error: ", err)
		return nil, err
	}
	return c.Do(args)
}

func (p *SSDBPool) routePool(args []string) (pool *gossdb.Connectors) {
	pool = p.MasterPool
	if len(args) < 2 {
		return
	}
	nk := args[1]
	for k, v := range p.PoolMap {
		if strings.HasPrefix(nk, k) {
			pool = v
			break
		}
	}
	return
}
