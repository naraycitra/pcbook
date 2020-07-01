gen:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:. --grpc-gateway_out=:. --openapiv2_out=logtostderr=true:./openapiv2
clean:
	rm pb/*.go
server:
	go run cmd/server/main.go -port 8083
rest:
	go run cmd/server/main.go -port 8084 -type rest -endpoint 0.0.0.0:8083
client:
	go run cmd/client/main.go -address 0.0.0.0:8081
test:
	go test -cover -race ./...
cert:
	cd certs; ./gen.sh; cd ..

.PHONY: gen clean server client test cert
