## Docker контейнеры для Expenses
При запуске 2-х и более контейнеров, для их взаимодействия дложна быть создана сеть.
```shell
docker network create expenses-net
```

### Контейнер expenses-pg (База данных)
> [!WARNING]
> При запуске приложения хост указанный в -h (например expenses-pg) должен совпадать с db_host указанным в config/main.yaml файле.
```shell
docker build -t expenses-pg --file ./build/docker/postgres.Dockerfile .
docker run -d --name expenses-pg \
    --env-file ./config/docker/.env.pg \
    -v expenses-pg:/var/lib/postgresql/data \
    --network expenses-net \
    -h expenses-pg \
    -p 5432:5432 expenses-pg
```

### Контейнер expenses-app (Приложение)
> [!ERROR]
> Не запустится без базы данных. Если нет expenses-pg, должен быть установлен PostgreSQL.
```shell
docker build -t expenses-pg --file ./build/docker/postgres.Dockerfile .
docker run -d --name expenses-pg \
    --env-file ./config/docker/.env.pg \
    -v expenses-pg:/var/lib/postgresql/data \
    --network expenses-net \
    -h expenses-pg \
    -p 5432:5432 expenses-pg
```

### Контейнеры expenses-app и expenses-pg (Приложение c базой данных)
Может пригодиться для разработки фронтенд части приложения. 
Но в таком случае, рекомендуется использовать `make start` или `docker compose up -d`
```shell
docker network create expenses-net

# Postgres
docker build -t expenses-pg --file ./build/docker/postgres.Dockerfile .
docker run -d --name expenses-pg \
    --env-file ./config/docker/.env.pg \
    -v expenses-pg:/var/lib/postgresql/data \
    --network expenses-net \
    -h expenses-pg \
    -p 5432:5432 expenses-pg

docker exec -it expenses-pg /bin/bash
psql -h expenses-pg -p 5432 -U username -W

# App
docker build -t expenses-app --file ./build/docker/app.Dockerfile .
docker run -d --name expenses-app \
    --network expenses-net \
    -p 7070:7070 expenses-app
```
