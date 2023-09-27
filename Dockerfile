FROM golang:1.20

RUN mkdir -p /app
ADD ./ /app/
WORKDIR /app

RUN go mod download && go mod verify
RUN go build -o bin ./cmd/app

EXPOSE 8080

CMD ["/app/bin"]