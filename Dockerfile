FROM golang:1.21.0-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/services/product_service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o salespot ./main.go

FROM alpine:3.18
WORKDIR /app
RUN chown nobody:nobody /app
USER nobody:nobody
COPY --from=builder --chown=nobody:nobody ./app/services/product_service/salespot .
COPY --from=builder --chown=nobody:nobody ./app/services/product_service/run.sh .

ENTRYPOINT sh run.sh
