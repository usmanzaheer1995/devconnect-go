package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/usmanzaheer1995/devconnect-go-v2/cmd/web/router"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/models/postgres"

	"github.com/usmanzaheer1995/devconnect-go-v2/cmd/web/config"
)

func main() {
	boolPtr := flag.Bool(
		"prod",
		false,
		"Provide this flag in production. This ensures that a config.json is provided before the application starts",
	)
	flag.Parse()

	cfg := config.LoadConfig(*boolPtr)
	services, err := postgres.NewServices(
		postgres.WithGorm(cfg.Database.ConnectionInfo()),
		postgres.WithUser(),
		postgres.WithProfile(),
	)

	if err != nil {
		panic(err)
	}
	r := router.NewRouter(services)

	defer services.Close()
	services.AutoMigrate()

	fmt.Printf("Starting server on port: %v\n", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r))
}
