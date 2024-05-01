FROM golang:1.22.2

WORKDIR /ipgem-api

RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY .air.toml .air.toml

CMD ["air", "-c", ".air.toml"]