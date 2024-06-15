```text
docker build -t expenses-pg:v.0 --file ./build/docker/postgres.Dockerfile .
docker volume create expenses-pg
docker run -d --name expenses-pg \
    -e POSTGRES_PASSWORD=postgres \
    -v expenses-pg:/var/lib/postgresql/data \
    -p 5432:5432 expenses-pg:v.0
docker exec -it expenses-pg /bin/bash
psql -h 0.0.0.0 -p 5432 -U postgres -W

create table users (
    id serial primary key,
    name varchar(64)
);

INSERT INTO users (name) VALUES ('Petr'), ('Ivan'), ('Boris');

select * from users;

go run cmd/expenses.go
```