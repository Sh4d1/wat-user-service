build:
	protoc -I/usr/local/include -I. \
		--go_out=plugins=micro:. \
		proto/user/user.proto
	docker build -t sh4d1/wat-user-service .

run:
	docker run --net="host" \
		-p 50052 \
		-e MICRO_ADDRESS=":50052" \
		-e MICRO_REGISTRY="mdns" \
		-e DB_HOST=localhost \
		-e DB_PASS=password \
		-e DB_USER=postgres \
		-e DEV=true \
		sh4d1/wat-user-service 
