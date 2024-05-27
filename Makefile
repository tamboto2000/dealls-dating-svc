build:
	@echo "building binary..."
	@go build -o bin/dealls-dating-svc cmd/main.go
	@echo "building finished"

local-deploy:
	docker compose -f ./deployment/local/docker-compose.yaml up -d

local-down:
	docker compose -f ./deployment/local/docker-compose.yaml down
