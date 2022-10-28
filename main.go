package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/vladoohr/simple_bank/api"
	db "github.com/vladoohr/simple_bank/db/sqlc"
	_ "github.com/vladoohr/simple_bank/doc/statik"
	"github.com/vladoohr/simple_bank/gapi"
	"github.com/vladoohr/simple_bank/pb"
	"github.com/vladoohr/simple_bank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load the configuration:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	// run migration
	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	go runGrpcServer(config, store)

	runGatewayServer(config, store)
}

func runDBMigration(url, dbSource string) {
	log.Println("Run DB migration")

	m, err := migrate.New(url, dbSource)
	if err != nil {
		log.Fatal("failed create new migrate instance:", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("failed run DB migrate up:", err)
	}

	log.Println("Migrations run successfully")
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create new grpc server: ", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	grpcMux := runtime.NewServeMux()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("failed to register handler server", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}

	log.Printf("start HTTP Gateway server on: %s", listener.Addr().String())

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("failed to start HTTP Gateway server: ", err)
	}
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create new grpc server: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("start gRPC server at: %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("failed to start gRPC server: ", err)
	}

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create new http server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start the server: ", err)
	}
}
