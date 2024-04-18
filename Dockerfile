FROM golang:1.20.13-alpine as rest_builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/rest ./cmd/rest/rest.go

FROM golang:1.20.13-alpine as migration

RUN apk add --no-cache make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

CMD make reset-db

FROM golang:1.20.13-alpine as rest_watch

RUN apk add --no-cache make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

CMD make reload

FROM alpine:3 as rest

WORKDIR /app

COPY --from=rest_builder /app/bin/rest /app/rest

EXPOSE 8080

CMD ./rest
