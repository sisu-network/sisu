FROM golang:1.18-alpine as builder

# This file is used in local dev with debugging purpose.

ENV GO111MODULE=on

WORKDIR /tmp/go-app

RUN apk add --no-cache make gcc musl-dev linux-headers git \
    && apk add openssh \
    && mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/sisu ./cmd/sisud/main.go

# Start fresh from a smaller image
FROM alpine:3.9

WORKDIR /app

COPY .env.dev /app/.env
COPY ./misc /app/misc
COPY --from=builder /tmp/go-app/out/sisu /app/sisu

RUN mkdir -p ~/.sisu

# Copy config into the container.
COPY docker-config/local-tss.toml /root/.sisu/tss/tss.toml
COPY docker-config/app-config.toml /root/.sisu/main/config/config.toml

CMD ["./sisu","start"]
