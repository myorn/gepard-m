package db

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/myorn/gepard-m/constants"
)

func InitDB() *sql.DB {
	// Connect to the database
	connStr := "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	return db
}

func DBUnaryServerInterceptor(session *sql.DB) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		return handler(context.WithValue(ctx, constants.DBSession, session), req)
	}
}
