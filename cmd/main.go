package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
	"github.com/lucidconnect/silver-arrow/erc20"
	"github.com/lucidconnect/silver-arrow/logger"
	"github.com/lucidconnect/silver-arrow/repository"
	"github.com/lucidconnect/silver-arrow/server"
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

	database := repository.NewPostgresDB(db)
	database.RunMigrations()
	// jobRunner := scheduler.NewScheduler(database)
	// setupJobs(jobRunner)

	httpServer := server.NewServer(database)
	httpServer.Routes()
	setupJobs(httpServer)
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	httpServer.Start(port)
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
	erc20.LoadSupportedChains("tokens/chains.json")
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

func setupJobs(runner *server.Server) {
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
