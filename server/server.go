package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	checkout_graph "github.com/lucidconnect/silver-arrow/graphql/checkout/graph"
	checkout_generated "github.com/lucidconnect/silver-arrow/graphql/checkout/graph/generated"
	merchant_graph "github.com/lucidconnect/silver-arrow/graphql/merchant/graph"
	merchant_generated "github.com/lucidconnect/silver-arrow/graphql/merchant/graph/generated"
	payment_link_graph "github.com/lucidconnect/silver-arrow/graphql/paymentLink/graph"
	payment_link_generated "github.com/lucidconnect/silver-arrow/graphql/paymentLink/graph/generated"
	wallet_graph "github.com/lucidconnect/silver-arrow/graphql/wallet/graph"
	wallet_generated "github.com/lucidconnect/silver-arrow/graphql/wallet/graph/generated"
	"github.com/lucidconnect/silver-arrow/repository"
	"github.com/lucidconnect/silver-arrow/service/erc4337"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"github.com/wader/gormstore/v2"
)

type Server struct {
	queue        repository.Queuer
	router       *mux.Router
	bundler      *erc4337.AlchemyService
	database     repository.Database
	sessionStore sessions.Store
	// walletGraphqlHandler, merchantGraphqlHandler *handler.Server
}

func NewServer(db *repository.PostgresDB) *Server {
	queue := repository.NewDeque()
	chain := os.Getenv("DEFAULT_CHAIN")
	defaultChain, err := strconv.ParseInt(chain, 10, 64)
	if err != nil {
		panic(err)
	}

	bundler, err := erc4337.NewAlchemyService(defaultChain)
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	sessionSecret := os.Getenv("JWT_SECRET")
	sesisonStore := gormstore.New(db.Db, []byte(sessionSecret))

	sesisonStore.SessionOpts = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	quit := make(chan struct{})
	go sesisonStore.PeriodicCleanup(1*time.Hour, quit)
	loadCORS(router)
	return &Server{
		queue:        queue,
		router:       router,
		bundler:      bundler,
		database:     db,
		sessionStore: sesisonStore,
	}
}

func (s *Server) Start(port string) {
	log.Info().Msgf("connect to http://localhost:%v/ for api/GraphQL playground", port)
	if err := http.ListenAndServe(":"+port, s.router); err != nil {
		log.Fatal().Err(err).Msg("unable to start the server")
	}
}

func (s *Server) Routes() {

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fmt.Sprintf("Lucid Backend Service %v",
			strings.ToTitle(os.Getenv("APP_ENV"))))
	})
	// merchant authentication
	authRouter := s.router.PathPrefix("/auth").Subrouter()
	nonceRoute := authRouter.HandleFunc("/nonce", s.GetNonce()).Methods(http.MethodGet)
	siweVerifyRoute := authRouter.HandleFunc("/verify", s.VerifyMerchant())
	authRouter.HandleFunc("/jwt", s.GenerateJwt()).Methods(http.MethodPost)
	authRouter.Use(MiddlewareExcept(s.BasicAuthMiddleware(), nonceRoute, siweVerifyRoute))

	merchantRouter := s.router.PathPrefix("/merchant").Subrouter()
	merchantRouter.Use(s.JWTMiddleware())
	merchantRouter.Handle("/graphiql", playground.Handler("api/GraphQL playground", "/merchant/query"))
	merchantRouter.Handle("/query", s.merchantGraphqlHandler())

	// checkout
	checkout := s.router.PathPrefix("/checkout").Subrouter()
	checkout.Use(s.CheckoutMiddleware())
	checkout.Handle("/query", s.checkoutGraphqlHandler())
	checkout.Handle("/graphiql", playground.Handler("/api/Graphql playground", "/checkout/query"))

	// wallet
	wallet := s.router.PathPrefix("/wallet").Subrouter()
	wallet.Handle("/query", s.walletGraphqlHandler())
	wallet.Handle("/graphiql", playground.Handler("/api/Graphql playground", "/wallet/query"))

	// payment page
	paymentLink := s.router.PathPrefix("/pay").Subrouter()
	paymentLink.Use(s.PaymentLinkMiddleware())
	paymentLink.Handle("/query", s.paymentLinksGraphqlHandler())
	paymentLink.Handle("/graphiql", playground.Handler("/api/Graphql plaground", "/pay/query"))
}

func (s *Server) checkoutGraphqlHandler() *handler.Server {
	checkoutsrv := handler.NewDefaultServer(checkout_generated.NewExecutableSchema(checkout_generated.Config{Resolvers: &checkout_graph.Resolver{
		Cache:    repository.NewMCache(),
		Database: s.database,
	}}))
	return checkoutsrv
}

func (s *Server) paymentLinksGraphqlHandler() *handler.Server {
	paymentLinkSrv := handler.NewDefaultServer(payment_link_generated.NewExecutableSchema(payment_link_generated.Config{Resolvers: &payment_link_graph.Resolver{
		Cache:    repository.NewMCache(),
		Database: s.database,
	}}))
	return paymentLinkSrv
}

func (s *Server) walletGraphqlHandler() *handler.Server {
	walletSrv := handler.NewDefaultServer(wallet_generated.NewExecutableSchema(wallet_generated.Config{Resolvers: &wallet_graph.Resolver{
		Cache:    repository.NewMCache(),
		Database: s.database,
	}}))
	return walletSrv
}

func (s *Server) merchantGraphqlHandler() *handler.Server {
	merchantSrv := handler.NewDefaultServer(merchant_generated.NewExecutableSchema(merchant_generated.Config{Resolvers: &merchant_graph.Resolver{
		Database: s.database,
	}}))
	return merchantSrv
}

func loadCORS(router *mux.Router) {
	switch os.Getenv("APP_ENV") {
	case "production":
		{
			allowedOrigins := []string{"https://portal.lucidconnect.xyz", "https://checkout.lucidconnect.xyz", "https://lucidconnect.xyz", "https://wallet.lucidconnect.xyz", "https://*", "http://*"}
			// for i := range utils.CustomMerchantCodes {
			// 	allowedOrigins = append(allowedOrigins, fmt.Sprintf("https://%v.web3-pay.com", utils.CustomMerchantCodes[i]))
			// }
			c := cors.New(cors.Options{
				AllowedOrigins: allowedOrigins,
				AllowedMethods: []string{
					http.MethodOptions,
					http.MethodGet,
					http.MethodPost,
				},
				AllowedHeaders:   []string{"*"},
				AllowCredentials: true,
			})
			c.Log = &log.Logger
			router.Use(c.Handler)
		}
	// case "staging":
	// 	c := cors.New(cors.Options{
	// 		// AllowedOrigins: []string{"https://checkout.sendcashpay.com", "https://*", "http://*", "https://checkout.transfers.africa"},
	// 		AllowedOrigins: []string{"https://portal.lucidconnect.xyz", "https://checkout.lucidconnect.xyz", "https://lucidconnect.xyz", "https://wallet.lucidconnect.xyz", "https://*", "http://*"},
	// 		AllowedMethods: []string{
	// 			http.MethodOptions,
	// 			http.MethodGet,
	// 			http.MethodPost,
	// 		},
	// 		AllowedHeaders:   []string{"*"},
	// 		AllowCredentials: true,
	// 	})
	// 	c.Log = &log.Logger
	// 	router.Use(c.Handler)
	default:
		c := cors.New(cors.Options{
			AllowedOrigins: []string{"https://portal.lucidconnect.xyz", "https://checkout.lucidconnect.xyz", "https://lucidconnect.xyz", "https://wallet.lucidconnect.xyz", "http://localhost:4002", "http://localhost:7890", "http://localhost:3000", "https://*"},
			AllowedMethods: []string{
				http.MethodOptions,
				http.MethodGet,
				http.MethodPost,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		})
		c.Log = &log.Logger
		router.Use(c.Handler)
	}
}

// func loadMerchantAuthMiddleware(router *mux.Router, db repository.Database) {
// 	// router.Use(auth.Middleware(db))
// 	router.Use()
// }
