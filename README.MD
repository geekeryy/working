## Working 项目管理工具

### 实现功能

> github webhook 触发k8s服务滚动更新

git提交如下格式commit，触发滚动更新

deploy:v{版本号}

版本号请参考[版本管理（ SemVer ）](https://github.com/comeonjy/go-layout#%E7%89%88%E6%9C%AC%E7%AE%A1%E7%90%86-semver-)

Example:
```shell
git commit --allow-empty -am "deploy:v0.0.1"
```

### TODO LIST

- [ ] 区分不同环境部署