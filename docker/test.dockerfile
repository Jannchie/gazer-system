FROM golang:alpine as builder

WORKDIR /app

RUN apk add build-base
COPY go.* .
RUN go mod download
COPY ../ ./

RUN --mount=type=cache,target=/root/.cache/go-build go build -o /out/client ./test/test.go

FROM alpine:latest as prod

COPY --from=builder /out/client /
CMD /client