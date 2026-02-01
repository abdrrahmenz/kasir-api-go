package database

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Add connection timeout and SSL parameters
	connStr := connectionString
	if len(connectionString) > 0 && !strings.Contains(connectionString, "sslmode=") {
		if strings.Contains(connectionString, "?") {
			connStr = connectionString + "&sslmode=require&connect_timeout=10"
		} else {
			connStr = connectionString + "?sslmode=require&connect_timeout=10"
		}
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Set connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test connection with timeout
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Database ping failed: %v", err)
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully")
	return db, nil
}
