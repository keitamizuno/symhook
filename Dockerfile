FROM golang:1.13 as builder
WORKDIR /go/src/app
COPY /src/ .
RUN go get github.com/dgrijalva/jwt-go && \
    go get github.com/gorilla/mux && \
    go get github.com/sirupsen/logrus
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/symhook .

FROM alpine:latest
COPY --from=builder /go/src/app/bin/symhook /app/
RUN mkdir -p /config
CMD ["/app/symhook"]