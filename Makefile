fe:
	go run cmd/frontend/frontend.go
be-auth:
	go run cmd/backend/backend.go -service=auth
be-adder:
	go run cmd/backend/backend.go -service=adder
be-multiply:
	go run cmd/backend/backend.go -service=multiply

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/*.proto


run-jaeger:
	docker run --name jaeger -p 13133:13133 -p 16686:16686 -p 4317:55680 -d --restart=unless-stopped jaegertracing/opentelemetry-all-in-one

.PHONY: proto