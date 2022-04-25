package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/aiit2022-pbl-okuhara/traning/micchie/sample-app/config"
)

func NewDB() (*sql.DB, error) {
	c := config.Config
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUsername, c.DBPassword, c.DBName)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return conn, nil
}
