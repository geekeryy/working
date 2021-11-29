
deploy:
	GOOS=linux GOARCH=amd64 go build -o main ./main.go
	docker build -t ccr.ccs.tencentyun.com/comeonjy/working:v0.0.1 .
	rm main
	echo "Aa123456" | docker login --username=1126254578 ccr.ccs.tencentyun.com --password-stdin
	docker push ccr.ccs.tencentyun.com/comeonjy/working:v0.0.1
	kubectl set image deploy working working=ccr.ccs.tencentyun.com/comeonjy/working:v0.0.1
	kubectl rollout restart deploy working
