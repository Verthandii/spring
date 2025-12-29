package ies

import (
	"context"
	"errors"
	"os"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v9"

	"github.com/Verthandii/spring/ilogger"
)

type ESClient struct {
	*elasticsearch.TypedClient
	Config *Config
}

// Init 初始化 ES TypedClient
func Init(cfg *Config) (*ESClient, error) {
	ilogger.Info("【es】 start init")
	clientCfg := elasticsearch.Config{
		Addresses:    []string{cfg.Address},
		Username:     cfg.Username,
		Password:     cfg.Password,
		DisableRetry: cfg.DisableRetry,
	}
	if cfg.EnableDebugLogger {
		clientCfg.EnableDebugLogger = cfg.EnableDebugLogger
		clientCfg.Logger = &elastictransport.JSONLogger{Output: os.Stdout, EnableRequestBody: true}
	}
	client, err := elasticsearch.NewTypedClient(clientCfg)
	if err != nil {
		return nil, err
	}
	success, err := client.Ping().Do(context.TODO())
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, errors.New("ping es error")
	}

	return &ESClient{Config: cfg, TypedClient: client}, nil
}
