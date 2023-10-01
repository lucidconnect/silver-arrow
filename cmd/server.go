package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"

	"github.com/helicarrierstudio/silver-arrow/graphql/wallet/graph"
	"github.com/helicarrierstudio/silver-arrow/graphql/wallet/graph/generated"
	"github.com/helicarrierstudio/silver-arrow/repository"
	"github.com/helicarrierstudio/silver-arrow/service/merchant"
	"github.com/helicarrierstudio/silver-arrow/service/scheduler"
	"github.com/helicarrierstudio/silver-arrow/service/turnkey"
	"github.com/helicarrierstudio/silver-arrow/service/wallet"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	loadEnv()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	db, err := repository.SetupDatabase(nil)
	if err != nil {
		log.Println(err)
	}
	router := chi.NewRouter()
	loadCORS(router)

	database := repository.NewDB(db)

	tunkeyService := turnkey.NewTurnKeyService()
	walletService := wallet.NewWalletService(database, tunkeyService)
	merchantService := merchant.NewMerchantService(database)

	router.Use(merchantService.Middleware())

	jobRunner := scheduler.NewScheduler(database, walletService)
	setupJobs(jobRunner)
	walletSrv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Database: database,
		Cache:    repository.NewMCache(),
	}}))

	merchantSrv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Database: database,
		Cache: repository.NewMCache(),
	}}))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/merchant/graphiql", playground.Handler("GraphQL playground", "/merchant/query"))

	router.Handle("/query", walletSrv)
	router.Handle("/merchant/query", merchantSrv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
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

func setupJobs(runner *scheduler.Scheduler) {
	log.Println("Setting up jobs...")
	c := cron.New(
		cron.WithChain(
			cron.Recover(cron.DefaultLogger),
			cron.SkipIfStillRunning(cron.DefaultLogger),
		),
	)

	c.AddFunc("@midnight", func() {
		runner.SubscriptionJob()
	})

	c.Start()
}

func loadCORS(router *chi.Mux) {
	switch os.Getenv("APP_ENV") {
	// case "production":
	// 	{
	// 		allowedOrigins := []string{"https://checkout.sendcashpay.com", "https://checkout.transfers.africa", "https://sendcashpay.com"}
	// 		for i := range utils.CustomMerchantCodes {
	// 			allowedOrigins = append(allowedOrigins, fmt.Sprintf("https://%v.web3-pay.com", utils.CustomMerchantCodes[i]))
	// 		}
	// 		router.Use(cors.New(cors.Options{
	// 			AllowedOrigins: allowedOrigins,
	// 			AllowedMethods: []string{
	// 				http.MethodOptions,
	// 				http.MethodGet,
	// 				http.MethodPost,
	// 			},
	// 			AllowedHeaders:   []string{"*"},
	// 			AllowCredentials: false,
	// 		}).Handler)
	// 	}
	// case "staging":
	// 	router.Use(cors.New(cors.Options{
	// 		AllowedOrigins: []string{"https://checkout.sendcashpay.com", "https://*", "http://*", "https://checkout.transfers.africa"},
	// 		AllowedMethods: []string{
	// 			http.MethodOptions,
	// 			http.MethodGet,
	// 			http.MethodPost,
	// 		},
	// 		AllowedHeaders:   []string{"*"},
	// 		AllowCredentials: false,
	// 	}).Handler)
	default:
		router.Use(cors.New(cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			AllowedMethods: []string{
				http.MethodOptions,
				http.MethodGet,
				http.MethodPost,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		}).Handler)
	}
}
