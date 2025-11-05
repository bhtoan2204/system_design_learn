# Hướng dẫn chạy Temporal Saga Orchestration

## Yêu cầu
1. Temporal server phải đang chạy (sử dụng docker-compose)
2. Worker phải được start trước khi chạy client

## Các bước chạy

### 1. Start Temporal Server
```bash
cd docker-compose
docker-compose up -d
```

### 2. Start Worker (Terminal 1)
Worker sẽ lắng nghe và xử lý workflows và activities:
```bash
cd orchestrator/worker
go run .
```

### 3. Start Client để trigger workflow (Terminal 2)
```bash
cd orchestrator/cmd
go run .
```

## Cấu trúc
- **Worker** (`orchestrator/worker/`): Xử lý workflows và activities
- **Client** (`orchestrator/cmd/`): Start workflow execution
- **Workflow**: Điều phối saga pattern với compensation logic
- **Activities**: Gọi HTTP API của các service

## Lưu ý
- Worker phải chạy trước khi client start workflow
- Các service HTTP endpoints (`/book`, `/cancel`) cần được implement
- Mặc định các service URLs:
  - Hotel: `http://localhost:8081`
  - Flight: `http://localhost:8082`
  - Car: `http://localhost:8083`

