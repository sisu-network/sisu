FROM golang:1.16-alpine as builder

ENV GO111MODULE=on \
    GOPRIVATE=github.com/sisu-network/*

WORKDIR /tmp/go-app

RUN apk add --no-cache make gcc musl-dev linux-headers git \
    && apk add openssh \
    && git config --global url."git@github.com:".insteadOf "https://github.com/" \
    && mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config

# # Though the id_rsa file is removed at the end of this docker build, it's still dangerous to include
# # id_rsa in the build file since docker build steps are cached. Only do this while our repos are in
# # private mode.
ADD tmp/id_rsa /root/.ssh/id_rsa

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/sisu ./cmd/sisud/main.go

RUN rm /root/.ssh/id_rsa

# Start fresh from a smaller image
FROM alpine:3.9 
RUN apk add ca-certificates

COPY --from=builder /tmp/go-app/out/sisu /app/sisu

CMD ["/app/sisu"]
