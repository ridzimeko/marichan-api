package main

import (
	"log"
	"marichan-api/internal/app"
	"marichan-api/internal/config"
)

func main() {
	env := config.LoadEnv()

	server, err := app.NewServer(env)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("server running on :%s\n", env.AppPort)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
