DOCKER_PSW:=Aa123456
DOCKER_USR:=1126254578
IMAGES_REPO:=ccr.ccs.tencentyun.com/comeonjy
REPO_DOMAIN:=ccr.ccs.tencentyun.com


CONTAINER_NAME:=working

# 重启服务
restart:
	kubectl set image deploy $(SERVER_NAME) $(CONTAINER_NAME)=$(IMAGES_REPO)/$(SERVER_NAME):$(IMAGE_TAG)
	kubectl rollout restart deploy $(SERVER_NAME)