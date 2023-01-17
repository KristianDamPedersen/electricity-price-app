FROM golang:latest

WORKDIR /app
ADD . .

RUN go vet ./...
RUN go build ./...


EXPOSE 8090

CMD ["./main", "serve", "--http=0.0.0.0:8090"]