FROM golang:1.16 as base
EXPOSE $PORT
WORKDIR /clean-boilerplate-go

FROM base AS development
RUN adduser gouser
USER gouser
RUN go get github.com/axw/gocov/gocov
RUN go get github.com/matm/gocov-html
RUN go get github.com/golang/mock/mockgen@v1.4.4

FROM base AS compiler-server
COPY . /clean-boilerplate-go
RUN go mod download
RUN CGO_ENABLED=0 GOARCH=amd64 go build -o /clean-boilerplate-go/application /clean-boilerplate-go/entrypoints/server/main.go

FROM alpine:3.10 AS base-prod
RUN apk add --update --no-cache ca-certificates tzdata && \
  rm -rf /var/cache/apk/* /tmp/* /var/tmp/* && \
  date

FROM base-prod AS production
COPY --from=compiler-server /clean-boilerplate-go/application /clean-boilerplate-go/application
RUN chmod +x /clean-boilerplate-go/application
CMD /clean-boilerplate-go/application
