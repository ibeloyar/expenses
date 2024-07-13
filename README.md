Getting started in Docker
```shell
docker network create expenses-net

# Postgres
docker build -t expenses-pg:0.0.1 --file ./build/docker/postgres.Dockerfile .
docker run -d --name expenses-pg \
    --env-file ./config/docker/.env.pg \
    -v expenses-pg:/var/lib/postgresql/data \
    --network expenses-net \
    -h expenses-pg \
    -p 5432:5432 expenses-pg:0.0.1

docker exec -it expenses-pg /bin/bash
psql -h expenses-pg -p 5432 -U username -W

# App
docker build -t expenses-app:0.0.1 --file ./build/docker/app.Dockerfile .
docker run -d --name expenses-app \
    --network expenses-net \
    -p 7070:7070 expenses-app:0.0.1
```

Make
```text
make run
make install-tools
make swagger-gen
make migrate-up
make migrate-down
make NAME="confirm_mail" migrate-create
```
