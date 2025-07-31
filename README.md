# Gin with Gorm 模板

### 项目组件
| 组件        | 说明       |
|-----------|----------|
| Web框架     | gin      |
| Console框架 | cobra    |
| Orm库      | gorm     |
| Redis库    | go-redis |
| 日志        | logrus   |
| 配置管理      | viper    |
| Beanstalkd队列 | go-beanstalk |

### 项目使用
+ 下载本项目后然后重命名包名
    ```bash
    ./rename-package.sh xxxx/xxx-api
    ```

### 开发热重启
+ 安装 `air-verse/air` 方便开发中热重启，无需放到 `go.mod` 中，然后直接运行 `air` 即可
    ```bash
    go install github.com/air-verse/air@v1.61.7

    air
    ```
