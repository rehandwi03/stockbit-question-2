ARG BASE_IMAGE=golang:1.16
FROM ${BASE_IMAGE}

RUN apt update -y && apt install upx -y

WORKDIR /src

COPY go.* ./
RUN go mod download

COPY . .

#RUN CGO_ENABLED=0 go build -ldflags="-s -w" -a && upx -q9 stockbit-service

RUN CGO_ENABLED=0 go build -o stockbit-service

FROM alpine:3.12

RUN apk add --update --no-cache tzdata ca-certificates

WORKDIR /app

COPY --from=0 /src/product-service .
COPY .env.example .env

EXPOSE 10000

ENTRYPOINT ["./stockbit-service"]
