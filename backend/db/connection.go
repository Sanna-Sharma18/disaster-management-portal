package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	go_ora "github.com/sijms/go-ora/v2"
)

func Connect() (*sql.DB, error) {
	host := getenv("ORACLE_HOST", "localhost")
	port := getenvInt("ORACLE_PORT", 1521)
	service := getenv("ORACLE_SERVICE", "XEPDB1")
	user := getenv("ORACLE_USER", "disaster_user")
	password := getenv("ORACLE_PASSWORD", "DisasterApp123!")

	connStr := go_ora.BuildUrl(host, port, service, user, password, nil)

	var (
		database *sql.DB
		err      error
	)

	// Retry up to 30 times (Oracle XE can take several minutes to start)
	for i := 1; i <= 30; i++ {
		database, err = sql.Open("oracle", connStr)
		if err == nil {
			err = database.Ping()
		}
		if err == nil {
			break
		}
		log.Printf("waiting for Oracle (%d/30): %v", i, err)
		time.Sleep(10 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("oracle connect: %w", err)
	}

	database.SetMaxOpenConns(20)
	database.SetMaxIdleConns(5)
	database.SetConnMaxLifetime(5 * time.Minute)
	return database, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
