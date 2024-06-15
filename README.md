```text
docker build -t expenses-pg:0.1.0 --file ./build/docker/postgres.Dockerfile .
docker volume create expenses-pg
docker run -d --name expenses-pg \
    -e POSTGRES_PASSWORD=postgres \
    -v expenses-pg:/var/lib/postgresql/data \
    -p 5432:5432 expenses-pg:0.1.0
docker exec -it expenses-pg /bin/bash
psql -h 0.0.0.0 -p 5432 -U postgres -W

docker build -t expenses-app:0.1.0 --file ./build/docker/app.Dockerfile .
docker run -d --name expenses-app -p 7070:7070 expenses-app:0.1.0
```