package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("Failed to parse database URL: %v", err)
	}

	// Configure connection pool
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	db, err := pgxpool.NewWithConfig(
		context.Background(),
		config,
	)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully")

	return db
}