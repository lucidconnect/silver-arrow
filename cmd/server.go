package main

import (
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
	"github.com/lucidconnect/silver-arrow/erc20"
	merchant_graph "github.com/lucidconnect/silver-arrow/graphql/merchant/graph"
	merchant_generated "github.com/lucidconnect/silver-arrow/graphql/merchant/graph/generated"
	wallet_graph "github.com/lucidconnect/silver-arrow/graphql/wallet/graph"
	wallet_generated "github.com/lucidconnect/silver-arrow/graphql/wallet/graph/generated"
	"github.com/lucidconnect/silver-arrow/logger"
	"github.com/lucidconnect/silver-arrow/repository"
	"github.com/lucidconnect/silver-arrow/service/merchant"
	"github.com/lucidconnect/silver-arrow/service/scheduler"
	"github.com/lucidconnect/silver-arrow/service/turnkey"
	"github.com/lucidconnect/silver-arrow/service/wallet"
	"github.com/robfig/cron/v3"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

const defaultPort = "8080"

var (
	db *gorm.DB
)

func main() {
	bootstrap()

	// db, err := repository.SetupDatabase(nil)
	// if err != nil {
	// 	log.Println(err)
	// }
	router := chi.NewRouter()
	loadCORS(router)

	database := repository.NewDB(db)
	database.RunMigrations()
	tunkeyService, err := turnkey.NewTurnKeyService()
	if err != nil {
		log.Panic().Err(err).Send()
	}
	walletService := wallet.NewWalletService(database, tunkeyService)
	merchantService := merchant.NewMerchantService(database)

	router.Use(merchantService.Middleware())

	jobRunner := scheduler.NewScheduler(database, walletService)
	setupJobs(jobRunner)
	walletSrv := handler.NewDefaultServer(wallet_generated.NewExecutableSchema(wallet_generated.Config{Resolvers: &wallet_graph.Resolver{
		Cache:          repository.NewMCache(),
		Database:       database,
		TurnkeyService: tunkeyService,
	}}))

	merchantSrv := handler.NewDefaultServer(merchant_generated.NewExecutableSchema(merchant_generated.Config{Resolvers: &merchant_graph.Resolver{
		Database: database,
	}}))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/merchant/graphiql", playground.Handler("GraphQL playground", "/merchant/query"))

	router.Handle("/query", walletSrv)
	router.Handle("/merchant/query", merchantSrv)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal().Err(err).Msg("unable to start the server")
	}
	log.Fatal()
}

func bootstrap() {
	var err error
	app := os.Getenv("APP_ENV")
	loadEnv(app)
	if app == "staging" || app == "production" {
		logger.SetUpLoggerFromConfig(app)
	} else {
		logger.SetUpDefaultLogger()
	}
	db, err = repository.SetupDatabase(nil)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to establish a database connection")
	}
	erc20.LoadSupportedTokens("tokens/tokens.json")
}

func loadEnv(app string) {
	switch app {
	case "development":
		log.Print("Loading configurations...Development")
		err := godotenv.Load(".env.development.local")
		if err != nil {
			log.Fatal().Err(err).Msg("Error loading .env file")
		}
	case "test":
		log.Print("Loading configurations...Test")
		err := godotenv.Load(".env.test.local")
		if err != nil {
			log.Fatal().Err(err).Msg("Error loading .env file")
		}
	}
}

func setupJobs(runner *scheduler.Scheduler) {
	log.Print("Setting up jobs...")
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
