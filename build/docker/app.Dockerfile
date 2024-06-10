FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download && go mod verify
RUN GOOS=linux GOARCH=amd64 go build -v -o expenses ./cmd/expenses.go


FROM alpine:3.20.0
WORKDIR /app
COPY --from=builder /app/expenses .
COPY --from=builder /app/config/main.yaml ./config/main.yaml
EXPOSE 7070
CMD ["./expenses"]