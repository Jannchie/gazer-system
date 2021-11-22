FROM golang:alpine as builder

WORKDIR /app

RUN apk add build-base
COPY go.* .
RUN go mod download
COPY ../ ./

RUN --mount=type=cache,target=/root/.cache/go-build go build -o /out/server ./cmd/server/main.go

FROM alpine:latest as prod

VOLUME /db
COPY --from=builder /out/server /
ENV GSS_PORT=2000 TOR="localhost:9050" TOR_CTL="localhost:9051" DSN="/db/data.db" TOR_PASSWORD="password" CONCURRENCY="8" SPEEDOS="http://localhost:4747"
CMD /server --port=${GSS_PORT} --tor=${TOR} --tor-ctl=${TOR_CTL} --dsn=${DSN} --tor-pass=${TOR_PASSWORD} --concurrency=${CONCURRENCY} --speedos=${SPEEDOS}
EXPOSE ${GSS_PORT}
