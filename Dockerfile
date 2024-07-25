FROM golang:1.22.5-alpine3.20 AS builder

#Копируем файлы с исходным кодом
COPY . /server

#Указываем рабочую директорию
WORKDIR /server

RUN go mod tidy
RUN go build -o ./bin/server ./cmd/main.go

FROM alpine:3.20.2

WORKDIR /root

COPY --from=builder /server/bin/server .

EXPOSE 8080

CMD ["./server"]