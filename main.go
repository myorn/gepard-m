package main

import (
	"database/sql"
	"log"
	"net"

	config "github.com/gookit/ini/v2"
	"google.golang.org/grpc"

	"github.com/myorn/gepard-m/constants"
	"github.com/myorn/gepard-m/cron"
	dblib "github.com/myorn/gepard-m/db"
	pb "github.com/myorn/gepard-m/dto/proto"
	"github.com/myorn/gepard-m/service"
)

func main() {
	loadConfig(constants.ConfigFilePath)

	dbSession := dblib.InitDB()
	dblib.Migrate(dbSession)

	defer dbSession.Close()

	serveGRPC(dbSession)

	go cron.RunCancelJob(dbSession)
}

func loadConfig(filename string) {
	// usually I prefer viper and yaml but it looks like an overkill here
	err := config.LoadFiles(filename)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}

func serveGRPC(dbSession *sql.DB) {
	port := config.String("server.port")

	// Create a gRPC server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			dblib.DBUnaryServerInterceptor(dbSession),
		),
	)

	pb.RegisterDepositServer(s, service.New())

	// Start the server
	log.Printf("Server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
