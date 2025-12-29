package source

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/constant"
	ac "github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/extension"
)

const (
	contentKey = "content"
)

func init() {
	extension.AddFormatParser(constant.JSON, &jsonExtParser{})
}

type apollo struct {
	client agollo.Client
	opt    *options
}

// Option is apollo option
type Option func(*options)

type options struct {
	appid          string
	secret         string
	cluster        string
	endpoint       string
	namespace      string
	isBackupConfig bool
	backupPath     string
}

// WithAppID with apollo config app id
func WithAppID(appID string) Option {
	return func(o *options) {
		o.appid = appID
	}
}

// WithCluster with apollo config cluster
func WithCluster(cluster string) Option {
	return func(o *options) {
		o.cluster = cluster
	}
}

// WithEndpoint with apollo config conf server ip
func WithEndpoint(endpoint string) Option {
	return func(o *options) {
		o.endpoint = endpoint
	}
}

// WithEnableBackup with apollo config enable backup config
func WithEnableBackup() Option {
	return func(o *options) {
		o.isBackupConfig = true
	}
}

// WithDisableBackup with apollo config enable backup config
func WithDisableBackup() Option {
	return func(o *options) {
		o.isBackupConfig = false
	}
}

// WithSecret with apollo config app secret
func WithSecret(secret string) Option {
	return func(o *options) {
		o.secret = secret
	}
}

// WithNamespace with apollo config namespace name
func WithNamespace(name string) Option {
	return func(o *options) {
		o.namespace = name
	}
}

// WithBackupPath with apollo config backupPath
func WithBackupPath(backupPath string) Option {
	return func(o *options) {
		o.backupPath = backupPath
	}
}

// NewApollo 创建一个Apollo数据源
func NewApollo(opts ...Option) Source {
	op := options{
		appid:          "name",
		cluster:        "env",
		endpoint:       "127.0.0.1",
		namespace:      "config.json",
		isBackupConfig: false,
	}
	for _, o := range opts {
		o(&op)
	}
	client, err := agollo.StartWithConfig(func() (*ac.AppConfig, error) {
		return &ac.AppConfig{
			AppID:            op.appid,
			Cluster:          op.cluster,
			NamespaceName:    op.namespace,
			IP:               op.endpoint,
			IsBackupConfig:   op.isBackupConfig,
			Secret:           op.secret,
			BackupConfigPath: op.backupPath,
		}, nil
	})
	if err != nil {
		panic(err)
	}
	return &apollo{client: client, opt: &op}
}

// Load 加载Apollo配置到数据源中
func (a *apollo) Load() ([]byte, error) {
	val, err := a.client.GetConfigCache(a.opt.namespace).Get(contentKey)
	if err != nil {
		return nil, err
	}
	return []byte(val.(string)), nil
}

type jsonExtParser struct{}

func (parser jsonExtParser) Parse(configContent any) (map[string]any, error) {
	return map[string]any{contentKey: configContent}, nil
}
