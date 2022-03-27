FROM golang:1.17 AS builder
COPY . /go/src/app
WORKDIR /go/src/app
RUN go mod download && go mod tidy -compat=1.17

#COPY build ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app cmd/main.go

FROM alpine
RUN apk add --no-cache tzdata
COPY --from=builder /app ./
EXPOSE $PORT
ENTRYPOINT ["./app"]
