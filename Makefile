IMAGE_TAG:=v0.0.1

SERVER_NAME:=working

# 临时
test:
	GOOS=linux GOARCH=amd64 go build -o main ./main.go
	docker build -t ccr.ccs.tencentyun.com/comeonjy/$(SERVER_NAME):$(IMAGE_TAG) .
	rm main
	echo "Aa123456" | docker login --username=1126254578 ccr.ccs.tencentyun.com --password-stdin
	docker push ccr.ccs.tencentyun.com/comeonjy/$(SERVER_NAME):$(IMAGE_TAG)
	kubectl set image deploy $(SERVER_NAME) working=ccr.ccs.tencentyun.com/comeonjy/$(SERVER_NAME):$(IMAGE_TAG)
	kubectl rollout restart deploy $(SERVER_NAME)

# 部署
deploy:
	GOOS=linux GOARCH=amd64 go build -o main ./main.go
	docker build -t ccr.ccs.tencentyun.com/comeonjy/$(SERVER_NAME):$(IMAGE_TAG) .
	rm main
	echo "Aa123456" | docker login --username=1126254578 ccr.ccs.tencentyun.com --password-stdin
	docker push ccr.ccs.tencentyun.com/comeonjy/$(SERVER_NAME):$(IMAGE_TAG)
	git commit --allow-empty -am "deploy:$(IMAGE_TAG)"
	git push

# 重启服务
restart:
	kubectl set image deploy $(SERVER_NAME) working=ccr.ccs.tencentyun.com/comeonjy/$(SERVER_NAME):$(IMAGE_TAG)
	kubectl rollout restart deploy $(SERVER_NAME)