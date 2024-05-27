package db

import (
    "database/sql"

    _ "github.com/lib/pq"
)

// ConnectDB connects to the PostgreSQL database
func ConnectDB() (*sql.DB, error) {
    connStr := "postgresql://neondb_owner:peE1ziNPkSx8@ep-divine-queen-a135fs1p-pooler.ap-southeast-1.aws.neon.tech/godb?sslmode=require"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    if err := db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}

// CreateUserSchema creates the user schema if it does not exist
func CreateUserSchema(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        password VARCHAR(100) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`
    _, err := db.Exec(query)
    return err
}
