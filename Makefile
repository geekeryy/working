# cat env.make
# DOCKER_PSW:=xxx
# DOCKER_USR:=xxx
# IMAGES_REPO:=ccr.ccs.tencentyun.com/xxx
# REPO_DOMAIN:=ccr.ccs.tencentyun.com
include env.make

# 镜像tag
IMAGE_TAG:=v0.0.1

SERVER_NAME:=working

# 自动生成文件
g:
	go generate -v .

# 代码检查
vet:
	 find * -type d -maxdepth 3 -print |  xargs -L 1  bash -c 'cd "$$0" && pwd  && go vet'

# 初始化
init:
	go env -w GO111MODULE=on
	go env -w GOPROXY=https://goproxy.cn,direct

# 部署
deploy:
	GOOS=linux GOARCH=amd64 go build -o main ./main.go
	docker build -t $(IMAGES_REPO)/$(SERVER_NAME):$(IMAGE_TAG) .
	rm main
	echo "$(DOCKER_PSW)" | docker login --username=$(DOCKER_USR) $(REPO_DOMAIN) --password-stdin
	docker push $(IMAGES_REPO)/$(SERVER_NAME):$(IMAGE_TAG)
	git commit --allow-empty -am "deploy:$(IMAGE_TAG)"
	git push



# 本地docker部署
docker:
	docker stop go-layout  & > /dev/null
	GOOS=linux GOARCH=amd64 go build -o main ./main.go
	docker build -t $(SERVER_NAME):$(IMAGE_TAG) .
	rm main
	docker run --rm -p 8080:8080 -p 8081:8081 -p 6060:6060 -d --name $(SERVER_NAME)  $(SERVER_NAME):$(IMAGE_TAG)
