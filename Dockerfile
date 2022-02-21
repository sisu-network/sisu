FROM golang:1.16-alpine as builder

# This file is used in local dev with debugging purpose.

ENV GO111MODULE=on \
    GOPRIVATE=github.com/sisu-network/*

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

#Workaround: We shouldn't make .env mandatory, and the environment variables can be loaded from multiple places.
# RUN apk add ca-certificates \
#     && touch /app/.env && echo "SAMPLE_KEY:SAMPLE_VALUE" > /app/.env

COPY .env.dev /app/.env
COPY tokens_dev.json /app/tokens_dev.json
COPY chains.json /app/chains.json
COPY --from=builder /tmp/go-app/out/sisu /app/sisu

RUN ./sisu localnet
RUN rm -rf ~/.sisu/main
RUN mv ./output/node0/main ~/.sisu/main

# Copy config into the container.
COPY docker-config/local-tss.toml /root/.sisu/tss/tss.toml
COPY docker-config/app-config.toml /root/.sisu/main/config/config.toml

CMD ["./sisu","start"]
