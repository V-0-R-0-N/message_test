package main

import (
	"context"
	"errors"
	"fmt"
	"message/internal/kafka"
	"message/internal/router"
	"message/internal/worker"
	storage "message/storage/postgres"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Функция для запуска сервера
func serverRun(s *http.Server) {
	// Получаем адрес и порт из переменных окружения
	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")
	addr := net.JoinHostPort(host, port)
	s.Addr = addr
	// Запускаем сервер
	fmt.Println("Server is running on :", addr)
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("ListenAndServe error: %v\n", err)
	}
}

func main() {

	// Создаем контекст, который будет отменен при получении сигналов os.Interrupt и syscall.SIGTERM
	// Нужно для правильной остановки внутри докера
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt,
		syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	// Инициализация хранилища
	db := storage.NewDB()
	defer db.Pool.Close()
	// Инициализация продюсера Kafka
	producer := kafka.InitProducer()
	defer (*producer).Close()

	// Инициализация роутера
	r := router.InitNew(db)

	// Создание структуры HTTP сервера
	server := &http.Server{
		Handler: r,
	}

	// Горутина для запуска сервера
	go serverRun(server)
	// Горутина для отправки сообщений в Kafka
	go worker.New(ctx, db, producer)

	// Ожидание сигнала о завершении
	<-ctx.Done()

	// Создание нового контекста с тайм-аутом для завершения сервера
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Остановка сервера
	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Server Shutdown Failed:%+v", err)
	} else {
		fmt.Println("Server exited properly")
	}
}
