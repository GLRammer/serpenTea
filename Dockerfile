FROM golang:1.26

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

EXPOSE 8080

COPY . .
RUN go build -o server .

CMD ["./server"]