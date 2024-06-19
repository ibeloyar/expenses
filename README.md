```text
docker build -t expenses-pg:v.0 --file ./build/docker/postgres.Dockerfile .
docker run -d --name expenses-pg \
    --env-file ./config/docker/.env.pg \
    -v expenses-pg:/var/lib/postgresql/data \
    -p 5432:5432 expenses-pg:v.0

docker exec -it expenses-pg /bin/bash
psql -h 0.0.0.0 -p 5432 -U postgres -W

docker build -t expenses-app:v.0 --file ./build/docker/app.Dockerfile .
docker run -d --name expenses-app -p 7070:7070 expenses-app:v.0
```