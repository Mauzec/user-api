package api

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/mauzec/user-api/db/sqlc"
	"github.com/mauzec/user-api/internal/token"
)

type TokenParams struct {
	AccessTokenDuration time.Duration
	// RefreshTokenDuration time.Duration
}

type Server struct {
	store  db.Store
	router *gin.Engine

	tokenMaker  token.Maker
	tokenParams TokenParams
}

func NewServer(store db.Store, tokenMaker token.Maker, tokenParams TokenParams) (*Server, error) {
	server := &Server{store: store, tokenMaker: tokenMaker, tokenParams: tokenParams}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("phone", validPhone); err != nil {
			return nil, fmt.Errorf("failed to register phone validation: %w", err)
		}
		if err := v.RegisterValidation("gender", validGender); err != nil {
			return nil, fmt.Errorf("failed to register gender validation: %w", err)
		}
	}

	server.SetupRouter()
	return server, nil
}

func (server *Server) SetupRouter() {
	router := gin.Default()

	_ = router.SetTrustedProxies(nil)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// single queries
	authRoutes.GET("/users/:username", server.getUserByUsername)
	authRoutes.POST("/users/:username", server.updateUser)

	server.router = router
}

func (server *Server) Run(addr string) error {
	return server.router.Run(addr)
}

// RunTLS starts the server with https suppor
func (server *Server) RunTLS(addr, certFile, keyFile string) error {
	return server.router.RunTLS(addr, certFile, keyFile)
}
