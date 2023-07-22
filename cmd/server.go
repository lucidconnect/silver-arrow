package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/helicarrierstudio/silver-arrow/graph"
	"github.com/helicarrierstudio/silver-arrow/graph/generated"
	"github.com/helicarrierstudio/silver-arrow/repository"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func main() {
	loadEnv()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mongoClient, err := repository.SetupMongoDatabase()
	if err != nil {
		log.Panic(err)
	}

	walletRepo := repository.NewMongoDb(mongoClient)
	walletSrv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		WalletRepository: walletRepo,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", walletSrv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func loadEnv() {
	env := os.Getenv("APP_ENV")

	switch env {
	case "development":
		log.Print("Loading configurations...Development")
		err := godotenv.Load(".env.development.local")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	case "test":
		log.Print("Loading configurations...Test")
		err := godotenv.Load(".env.test.local")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}
