FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o app .

CMD ["/app/app"]