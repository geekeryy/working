# cat env.make
# DOCKER_PSW:=xxx
# DOCKER_USR:=xxx
# IMAGES_REPO:=ccr.ccs.tencentyun.com/xxx
# REPO_DOMAIN:=ccr.ccs.tencentyun.com
include env.make

# 镜像tag
IMAGE_TAG:=v0.0.1

SERVER_NAME:=working

# 部署
deploy:
	GOOS=linux GOARCH=amd64 go build -o main ./main.go
	docker build -t $(IMAGES_REPO)/$(SERVER_NAME):$(IMAGE_TAG) .
	rm main
	echo "$(DOCKER_PSW)" | docker login --username=$(DOCKER_USR) $(REPO_DOMAIN) --password-stdin
	docker push $(IMAGES_REPO)/$(SERVER_NAME):$(IMAGE_TAG)
	git commit --allow-empty -am "deploy:$(IMAGE_TAG)"
	git push
