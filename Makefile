# 镜像仓库配置
# export DOCKER_PSW=xxx
# export DOCKER_USR=xxx
# export IMAGES_REPO=ccr.ccs.tencentyun.com/xxx
# export REPO_DOMAIN=ccr.ccs.tencentyun.com


# 镜像tag
IMAGE_TAG:=v0.0.1

SERVER_NAME:=working
CONTAINER_NAME:=working

# 自动生成文件
g:
	go generate -v .

# 代码检查
vet:
	 find * -type d -maxdepth 3 -print |  xargs -L 1  bash -c 'cd "$$0" && pwd  && go vet'

# 初始化
init:
	go env  GO111MODULE=on
	go env  GOPROXY=https://goproxy.cn,direct

# 部署到k8s
deploy:
	GOOS=linux GOARCH=amd64 go build -o main ./main.go
	docker build -t $(IMAGES_REPO)/$(SERVER_NAME):$(IMAGE_TAG) .
	rm main
	echo "$(DOCKER_PSW)" | docker login --username=$(DOCKER_USR) $(REPO_DOMAIN) --password-stdin
	docker push $(IMAGES_REPO)/$(SERVER_NAME):$(IMAGE_TAG)
	git commit --allow-empty -am "deploy:$(IMAGE_TAG)"
	git push

# 部署到k8s 不提交
update:
	GOOS=linux GOARCH=amd64 go build -o main ./main.go
	docker build -t $(IMAGES_REPO)/$(SERVER_NAME):$(IMAGE_TAG) .
	rm main
	echo "$(DOCKER_PSW)" | docker login --username=$(DOCKER_USR) $(REPO_DOMAIN) --password-stdin
	docker push $(IMAGES_REPO)/$(SERVER_NAME):$(IMAGE_TAG)
	kubectl rollout restart deploy $(SERVER_NAME)

# 本地docker部署
docker:
	docker stop $(SERVER_NAME)  & > /dev/null
	GOOS=linux GOARCH=amd64 go build -o main ./main.go
	docker build -t $(SERVER_NAME):$(IMAGE_TAG) .
	rm main
	docker run --rm -p 8080:8080 -p 8081:8081 -p 6060:6060 -d --name $(SERVER_NAME)  $(SERVER_NAME):$(IMAGE_TAG)

# dlv 远程k8s调试
k8s-dlv:
	GOOS=linux GOARCH=amd64 go build -gcflags="all=-N -l" -o main ./main.go
	docker build -t $(IMAGES_REPO)/$(SERVER_NAME):dlv -f dlv.Dockerfile .
	rm main
	echo "$(DOCKER_PSW)" | docker login --username=$(DOCKER_USR) $(REPO_DOMAIN) --password-stdin
	docker push $(IMAGES_REPO)/$(SERVER_NAME):dlv
	kubectl rollout restart deploy $(SERVER_NAME)-dlv

# 本地调试
debug-dev:export APP_ENV=dev
debug-dev:
	go build -gcflags "all=-N -l" main.go
	dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./main

# 更新依赖
u:
	go get github.com/comeonjy/go-kit@main

# 拦截k8s流量到本地
intercept:
	telepresence quit -u && telepresence connect
	telepresence intercept $(SERVER_NAME) --http-header=who=jiangyang -n default