build:
	sudo docker build -t go-api-service .
run:
	sudo docker run -d -v $(shell pwd)/.env:/app/.env -v $(shell pwd)/runtime:/app/runtime -p 8089:8089 go-api-service
clear:
	sudo docker rm $(sudo docker ps -a -q)