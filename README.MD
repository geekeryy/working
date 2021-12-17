## Working 项目管理工具

### 说明
此服务为快速构建部署提供支持

### 实现功能

> #### github webhook 触发k8s服务滚动更新

原理：获取git push事件，通过k8s api更新服务

1. github webhook 填写接口地址 例如：https://xxx.abc.com/github-event

2. commit格式：deploy:v{版本号}

    版本号请参考[版本管理（ SemVer ）](https://github.com/comeonjy/working#%E7%89%88%E6%9C%AC%E7%AE%A1%E7%90%86-semver-)

    Example:
```shell
git commit --allow-empty -am "deploy:v0.0.1"
```
3. deployment.yaml 添加labels: githubRepoName={github仓库名}

    Example:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: working
  namespace: default
  labels:
    app: working
    githubRepoName: working
    ...
```


> #### Apollo webhook 主动触发配置热更新

原理：通过Apollo webhook接收配置更新事件，从k8s api读取apollo关联服务podIP集合，最后广播到每个服务

1. 注册ReloadConfigServer（注册了ReloadConfigServer的服务才能主动触发热更新）

```go
import "github.com/comeonjy/go-kit/grpc/reloadconfig"
    
reloadconfig.RegisterReloadConfigServer(server, reloadconfig.NewServer(conf))
```

2. deployment.yaml 添加labels: apolloAppId={Apollo配置中心的AppId}

    Example:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: working
  namespace: default
  labels:
    app: working
    apolloAppId: working
    ...
```

### TODO LIST

- [ ] 区分不同环境部署