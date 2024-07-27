package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"message/models"
	"time"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB() *DB {
	// "postgres://username:password@localhost:5432/dbname"
	connString := "postgres://us:pass@postgres:5432/messaggio"
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

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal("Ping error")
	}

	return &DB{
		Pool: pool,
	}
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

func (db *DB) GetStats() *models.Stats {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resultSend := 0
	err := db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM messaggio.public.messages WHERE sent=false").Scan(&resultSend)
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
		return nil
	}

	return &models.Stats{
		Counter: resultSend,
	}
}
