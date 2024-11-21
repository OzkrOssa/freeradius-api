FROM golang:1.22-alpine as builder

WORKDIR /app
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/freeradius_api ./cmd/main.go

FROM alpine:latest AS final

WORKDIR /app

COPY --from=build /app/bin/freeradius_api ./

EXPOSE 8080

ENTRYPOINT [ "./freeradius_api" ]