FROM golang:1.19.5-alpine3.17
WORKDIR /app
ADD . .

RUN apk add build-base
RUN go vet ./...
RUN go build ./...


EXPOSE 8090

CMD ["go", "run", "./...", "serve", "--http=0.0.0.0:8090"]
