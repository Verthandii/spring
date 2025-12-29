package env

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/Verthandii/spring/db/icache"
	"github.com/Verthandii/spring/utils"
)

const (
	Local      = "LOCAL"      // 本地环境
	Dev        = "DEV"        // 开发环境
	Test       = "TEST"       // 测试环境
	Staging    = "STAGING"    // 预发布环境
	Production = "PRODUCTION" // 生产环境

	// AppEnv 应用程序环境
	AppEnv = "APP_ENV"
	// AppVersion 应用程序版本
	AppVersion = "APP_VERSION"
	// AppServer 应用程序所属服务
	AppServer = "APP_SERVER"
	// AppBuildAt 构建时间
	AppBuildAt = "APP_BUILD_AT"
	// AppRegion 应用程序区域
	AppRegion = "APP_REGION"
	// AppRootPath 应用程序根目录
	AppRootPath = "APP_ROOT_PATH"
	// AppDebug 设置为 1 表示开启 pprof 监控
	AppDebug = "APP_DEBUG"
	// AppDisableRequestLog 设置为 1 表示禁用请求日志
	AppDisableRequestLog = "APP_DISABLE_REQUEST_LOG"
	// AppLogLevel 日志等级 debug info warn error
	AppLogLevel = "APP_LOG_LEVEL"

	// AppApolloEndpoint 应用访问Apollo域名
	AppApolloEndpoint = "APP_APOLLO_ENDPOINT"
)

var (
	osenv     = icache.NewMCache()
	errAppEnv = errors.New("环境变量 APP_ENV 填写错误，请检查后重启服务")
)

func init() {
	wd, _ := os.Getwd()
	_ = Setenv(AppRootPath, wd)
}

// Setenv 设置环境变量
func Setenv(k, v string) error {
	osenv.Set(k, v)
	return nil
}

// getEnv 根据key获取环境变量
func getEnv(key string) string {
	if env, ok := osenv.Get(key); ok {
		return env.(string)
	}
	env := os.Getenv(key)
	_ = Setenv(key, env)
	return env
}

// GetEnv 获取 APP_ENV
func GetEnv() string {
	return strings.ToUpper(getEnv(AppEnv))
}

// GetLEnv 获取 小写 APP_ENV
func GetLEnv() string {
	return strings.ToLower(getEnv(AppEnv))
}

// GetBuildAt 获取 APP_BUILD_AT
func GetBuildAt() int64 {
	if env, ok := osenv.Get(AppBuildAt); ok {
		return env.(int64)
	}

	buildAt, _ := strconv.ParseInt(os.Getenv(AppBuildAt), 10, 64)
	osenv.Set(AppBuildAt, buildAt)

	return buildAt
}

// GetRootPath 获取 APP_ROOT_PATH
func GetRootPath() string {
	return getEnv(AppRootPath)
}

// IsLocal 是否本地环境
func IsLocal() bool {
	return GetEnv() == Local
}

// IsDev 是否开发环境
func IsDev() bool {
	return GetEnv() == Dev
}

// IsTest 是否测试环境
func IsTest() bool {
	return GetEnv() == Test
}

// IsStaging 是否预发布环境
func IsStaging() bool {
	return GetEnv() == Staging
}

// IsProd 是否正式环境
func IsProd() bool {
	return GetEnv() == Production
}

// IsStagingOrProd 是否预发布/生产环境
func IsStagingOrProd() bool {
	env := GetEnv()
	return Staging == env || Production == env
}

// GetURegion 获取大写的 APP_REGION
func GetURegion() string {
	return strings.ToUpper(getEnv(AppRegion))
}

// GetLRegion 获取小写的 APP_REGION
func GetLRegion() string {
	return strings.ToLower(getEnv(AppRegion))
}

// GetVersion 获取 APP_VERSION
func GetVersion() string {
	return getEnv(AppVersion)
}

// GetServer 获取 APP_SERVER
func GetServer() string {
	return getEnv(AppServer)
}

// ServerIs 是否为某个服务
func ServerIs(name string) bool {
	return GetServer() == name
}

// CheckEnv 检查环境变量
func CheckEnv() error {
	env := GetEnv()
	if !utils.InSlice(env, []string{Local, Dev, Test, Staging, Production}) {
		return errAppEnv
	}
	return nil
}

// DisableReqLog 是否禁用了请求日志
func DisableReqLog() bool {
	return getEnv(AppDisableRequestLog) == "1"
}

// EnableDebug 是否开启了调试模式【pprof】
func EnableDebug() bool {
	return getEnv(AppDebug) == "1"
}

// LogLevel 获取 APP_LOG_LEVEL
func LogLevel() string {
	return getEnv(AppLogLevel)
}

// GetApolloEndpoint 获取国服Apollo域名
func GetApolloEndpoint() string {
	host := getEnv(AppApolloEndpoint)
	if host == "" {
		return "http://apollo-inner.yostar.net:8080"
	}
	return host
}

// GetIntlApolloEndpoint 获取海外Apollo域名
func GetIntlApolloEndpoint() string {
	host := getEnv(AppApolloEndpoint)
	if host == "" {
		return "http://apollo-oversea-inter.yo-star.com:8080"
	}
	return host
}
