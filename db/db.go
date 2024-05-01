package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	config "github.com/gookit/ini/v2"
	"google.golang.org/grpc"

	"github.com/myorn/gepard-m/constants"
)

func InitDB() *sql.DB {
	// Connect to the database
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		config.String("db.user"), config.String("db.pw"), config.String("db.name"))
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
