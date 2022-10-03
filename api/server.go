package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/vladoohr/simple_bank/db/sqlc"
	"github.com/vladoohr/simple_bank/token"
	"github.com/vladoohr/simple_bank/util"
)

// Server server http requests
type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

// NewServer creates new http sesrver and setup routing
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", currencyValidator)
	}

	server.setUpRouter()

	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	router.POST("/users", server.CreateUser)
	router.POST("/users/login", server.LoginUser)
	router.POST("/tokens/renew_access", server.RenewAccessToken)

	authRoutes := router.Group("/").Use(AuthMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.CreateAccount)
	authRoutes.GET("/accounts/:id", server.GetAccount)
	authRoutes.GET("/accounts", server.ListAccount)

	authRoutes.POST("/transfers", server.CreateTransfer)

	server.router = router
}

// Start runs http server on  a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse returns error formatted for the HTTP response
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
