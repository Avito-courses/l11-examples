
Сборка образа. Проверка dockerfile
```bash
docker build -t my-app:local .
```

Запустить форматеры и линтеры
```bash
golangci-lint fmt
golangci-lint run
```

Поднять БД локально
```docker-compose up -d```

Пинг http://127.0.0.1:3005/ping

Запустить миграции
```bash
set -a
source .env
set +a

goose -dir=migrations postgres "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB" up
```
