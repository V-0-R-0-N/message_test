package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"message/models"
	"os"
	"strconv"
	"time"
)

// DB структура для работы с базой данных
type DB struct {
	Pool *pgxpool.Pool
}

// NewDB создает новый экземпляр DB
func NewDB() *DB {
	// Получение переменных окружения
	dbName := os.Getenv("SERVER_POSTGRES_DB")
	postgresUser := os.Getenv("SERVER_POSTGRES_USER")
	postgresPassword := os.Getenv("SERVER_POSTGRES_PASSWORD")
	timeout, err := strconv.Atoi(os.Getenv("SERVER_POSTGRES_TIMEOUT"))
	if err != nil {
		log.Fatalf("Failed to convert WORKER_TIMEOUT to int: %v", err)
	}
	connString := fmt.Sprintf("postgres://%s:%s@postgres:5432/%s", postgresUser, postgresPassword, dbName)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v\n", err)
	}
	config.MaxConns = 10
	config.MinConns = 1

	// Создание пула соединений
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	// Проверка соединения
	time.Sleep(time.Duration(timeout) * time.Second)
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal("Ping error")
	}

	return &DB{
		Pool: pool,
	}
	// TODO defer после вызова функции
}

func (db *DB) Save(req *models.Message) error {

	req.Created = time.Now()

	query := `
		INSERT INTO messaggio.public.messages (author, text, created, sent)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	ctx := context.Background()
	err := db.Pool.QueryRow(ctx, query, req.Author, req.Text, req.Created, req.Sent).Scan(&req.ID)
	if err != nil {
		log.Fatal("Unable to fetch message count")
		return err
	}
	return nil
}

// GetStats возвращает статистику по принятым и отправленным сообщениям
func (db *DB) GetStats() *models.Stats {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	totalSent := 0
	resultSend := 0
	err := db.Pool.QueryRow(ctx, "SELECT COUNT(id) FROM messaggio.public.messages").Scan(&totalSent)
	if err != nil {
		log.Printf("QueryRow total failed: %v\n", err)
		return nil
	}
	err = db.Pool.QueryRow(ctx, "SELECT COUNT(id) FROM messaggio.public.messages WHERE sent=true").Scan(&resultSend)
	if err != nil {
		log.Printf("QueryRow sent failed: %v\n", err)
		return nil
	}

	return &models.Stats{
		Total: totalSent,
		Sent:  resultSend,
	}
}

// NeedSent возвращает сообщения, которые необходимо отправить
func (db *DB) NeedSent() []*models.Message {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.Pool.Query(ctx, "SELECT id, author, text, created FROM messaggio.public.messages WHERE sent=false")
	if err != nil {
		log.Printf("Query failed: %v\n", err)
		return nil
	}
	defer rows.Close()

	var result []*models.Message
	for rows.Next() {
		var message models.Message
		err = rows.Scan(&message.ID, &message.Author, &message.Text, &message.Created)
		if err != nil {
			log.Printf("Scan failed: %v\n", err)
			return nil
		}
		result = append(result, &message)
	}
	return result
}

// ChangeStatusSent изменяет статус сообщения на отправленное
func (db *DB) ChangeStatusSent(ctx context.Context, id int) error {
	_, err := db.Pool.Exec(ctx, "UPDATE messaggio.public.messages SET sent=true WHERE id=$1", id)
	if err != nil {
		log.Printf("Exec failed: %v\n", err)
		return err
	}
	return nil
}
