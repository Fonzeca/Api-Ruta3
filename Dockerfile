# syntax=docker/dockerfile:1

FROM golang:alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /api-ruta3 ./src/

EXPOSE 2252

CMD [ "/api-ruta3" ]