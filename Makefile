.PHONY: dev seed test build clean docker-up docker-down

# 開發模式：啟動 Go API + Vite 前端（需先啟動 MongoDB）
dev:
	@echo "啟動 API 伺服器..."
	go run cmd/api/main.go &
	@echo "啟動前端開發伺服器..."
	cd web && npm run dev

# 種子資料
seed:
	go run cmd/seed-missions/main.go

# 測試
test:
	go test ./...

# 前端型別檢查
typecheck:
	cd web && npx tsc --noEmit

# 建置
build:
	go build -o bin/api ./cmd/api
	go build -o bin/seed ./cmd/seed-missions
	cd web && npm run build

# Docker Compose 啟動
docker-up:
	docker compose up -d
	@echo "等待 MongoDB 啟動..."
	sleep 3
	docker compose exec api seed
	@echo "✅ 服務已啟動：http://localhost:5173"

# Docker Compose 關閉
docker-down:
	docker compose down

# 清理
clean:
	rm -rf bin/ web/dist/
