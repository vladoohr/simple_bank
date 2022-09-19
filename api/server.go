package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/vladoohr/simple_bank/db/sqlc"
)

// Server server http requests
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates new http sesrver and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", currencyValidator)
	}

	router.POST("/users", server.CreateUser)

	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.GetAccount)
	router.GET("/accounts", server.ListAccount)

	router.POST("/transfers", server.CreateTransfer)

	server.router = router

	return server
}

// Start runs http server on  a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse returns error formatted for the HTTP response
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
