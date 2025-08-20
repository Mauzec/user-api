package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/mauzec/user-api/db/sqlc"
	"github.com/mauzec/user-api/internal/api"
	"github.com/mauzec/user-api/internal/config"
	"github.com/mauzec/user-api/internal/token"
)

func main() {
	config, err := config.LoadConfig("app", "env", "./config")
	if err != nil {
		log.Fatal("unable to load config:", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("unable to connect to database:", err)
	}
	defer conn.Close()

	var tokenMaker token.Maker
	switch config.TokenType {
	case "PasetoS":
		tokenMaker, err = token.NewPasetoSMaker(config.TokenSymmetricKey)
	default:
		log.Fatal("given unsupported token type")
	}
	if err != nil {
		log.Fatal("something go wrong when creating token maker:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store, tokenMaker, api.TokenParams{
		AccessTokenDuration: config.AccessTokenDuration,
	})
	if err != nil {
		log.Fatal("server creating err:", err)
	}

	if config.TLSCertFile != "" && config.TLSKeyFile != "" {
		err = server.RunTLS(config.ServerAddr, config.TLSCertFile, config.TLSKeyFile)
	} else {
		err = server.Run(config.ServerAddr)
	}
	if err != nil {
		log.Fatal("catch error when starting server")
	}
}
