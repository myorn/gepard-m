package db

import (
	"context"
	"database/sql"
	"log"
)

const initTables = `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Create the messages table
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    message JSONB,
    tx_id BIGINT UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
-- Create the deposit table
CREATE TABLE IF NOT EXISTS deposit (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    credit BIGINT DEFAULT 0 CHECK (credit >= 0),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO deposit(credit) VALUES(0);`

const duplicateTableErrorCode = "42P07"

type checker interface {
	SQLState() string
}

func Migrate(db *sql.DB) {
	_, err := db.ExecContext(context.Background(), initTables)
	if err != nil {
		pe, ok := err.(checker)
		// not a postgre error or not a duplicate error
		if !ok || pe.SQLState() != duplicateTableErrorCode {
			log.Fatalf("Failed to migrate, %v", err)
		}

		log.Printf("Already migrated, %v", err)
	}
}

func Flush(db *sql.DB) {
	_, err := db.ExecContext(context.Background(), `TRUNCATE messages, deposit RESTART IDENTITY`)
	if err != nil {
		log.Fatalf("Failed to flush shema, %v", err)
	}
}
