FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o app . && apt update && apt-get install -y iputils-ping 

CMD ["/app/app"]