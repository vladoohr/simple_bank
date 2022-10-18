package gapi

import (
	"fmt"

	db "github.com/vladoohr/simple_bank/db/sqlc"
	"github.com/vladoohr/simple_bank/pb"
	"github.com/vladoohr/simple_bank/token"
	"github.com/vladoohr/simple_bank/util"
)

// Server server grpc requests
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	pb.UnimplementedSimpleBankServer
}

// NewServer creates new grpc sesrver and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker([]byte(config.TokenSymmetricKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
