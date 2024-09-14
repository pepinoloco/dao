package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"dao/cmd/user-server/internal/database"
	"dao/cmd/user-server/internal/env"
	"dao/cmd/user-server/internal/httpx"
	"dao/pkg/api/user"
	"dao/pkg/dal"
	"github.com/google/uuid"
)

var info = httpx.Info{
	Name:    "user-dal-example",
	Version: "v0.1.0",
}

func main() {
	ctx := context.Background()

	config := env.Load()

	db, err := database.Open(ctx, []string{dal.UserTable})
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	u := &user.Handler{
		UserDAO: &dal.User{
			DB: db,
			GenerateUUID: func() string {
				return uuid.New().String()
			},
		},
	}

	server := httpx.NewServer(info, config.HTTP.Port, config.HTTP.ReadTimeout, config.HTTP.WriteTimeout)
	server.Start([]httpx.Router{u})
	defer func() {
		if err := server.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	log.Printf("server running on port %s", config.HTTP.Port)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	_ = <-sig
}
