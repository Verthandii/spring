# 环境变量

| 名称                      | 值                                        | 来源                          |
|-------------------------|------------------------------------------|-----------------------------| 
| APP_ENV                 | 不区分大小写 LOCAL DEV TEST STAGING PRODUCTION | k8s deployment => env       |
| APP_VERSION             | x.y.z.n                                  | devops => Dockerfile => env |
| APP_BUILD_AT            | 秒级时间戳                                    | devops => Dockerfile => env |
| APP_REGION              | 不区分大小写 cn jp us                          | k8s deployment => env       |
| APP_SERVER              | TODO                                     | k8s deployment => env       |
| APP_ROOT_PATH           | 程序启动目录                                   | program => env              |
| APP_DEBUG               | 设置为 1 表示开启 pprof 监控                      | k8s deployment => env       |
| APP_DISABLE_REQUEST_LOG | 设置为 1 表示禁用请求日志                           | k8s deployment => env       |
| APP_LOG_LEVEL           | 日志等级 DEBUG INFO WARN ERROR               | k8s deployment => env       |
