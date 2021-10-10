FROM golang:alpine as builder

WORKDIR /app

RUN apk add build-base
COPY go.* .
RUN go mod download
COPY ../ ./

RUN --mount=type=cache,target=/root/.cache/go-build go build -o /out/server ./cmd/gs/main.go

FROM alpine:latest as prod

VOLUME /db
COPY --from=builder /out/server /
CMD /server

EXPOSE 2000
