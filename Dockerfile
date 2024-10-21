FROM golang:1.22-alpine as build

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/api ./cmd/main.go

FROM alpine:latest AS final

WORKDIR /app

COPY --from=build /app/bin/api ./

EXPOSE 8000

CMD ["./api"]
