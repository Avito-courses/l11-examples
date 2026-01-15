## Github actions
Конфиг в `.github/workflows/ci.yml`  
Собранные образы можно увидеть [тут](https://github.com/Avito-courses/l11-examples/pkgs/container/l11-examples) 

## Линтеры
Конфиг линтеров лежит в `.golangci.yaml`  
Запустить форматеры и линтеры
```bash
golangci-lint fmt
golangci-lint run
```

## Dockerfile
Сборка образа. Проверка dockerfile
```bash
docker build -t my-app:local .
```

---

Поднять БД локально
```docker-compose up -d```

Запустить миграции
```bash
set -a
source .env
set +a

goose -dir=migrations postgres "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB" up
```
