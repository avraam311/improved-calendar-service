FROM golang:alpine AS build_base

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN go build -o app ./cmd/main.go

FROM alpine AS runner

COPY --from=build_base /app .
COPY ./.env .
COPY ./config/config.yaml ./config/config.yaml

CMD ["./app"]