FROM golang:1.22.5-alpine3.20 AS builder

#Копируем файлы с исходным кодом
COPY . /server

#Указываем рабочую директорию
WORKDIR /server

#Установка записимостей и сборка бинарника
RUN go mod tidy
RUN go build -o ./bin/server ./cmd/main.go

FROM alpine:3.20.2

#Указываем рабочую директорию
WORKDIR /root

#Копируем бинарь из контейнера builder
COPY --from=builder /server/bin/server .

CMD ["./server"]