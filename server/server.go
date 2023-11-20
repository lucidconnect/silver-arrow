package server

import (
	"net/http"
	"os"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/lucidconnect/silver-arrow/auth"
	merchant_graph "github.com/lucidconnect/silver-arrow/graphql/merchant/graph"
	merchant_generated "github.com/lucidconnect/silver-arrow/graphql/merchant/graph/generated"
	wallet_graph "github.com/lucidconnect/silver-arrow/graphql/wallet/graph"
	wallet_generated "github.com/lucidconnect/silver-arrow/graphql/wallet/graph/generated"
	"github.com/lucidconnect/silver-arrow/repository"
	"github.com/lucidconnect/silver-arrow/service/erc4337"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

type Server struct {
	queue        repository.Queuer
	router       *mux.Router
	bundler      *erc4337.AlchemyService
	database     repository.Database
	sessionStore *sessions.CookieStore
	// walletGraphqlHandler, merchantGraphqlHandler *handler.Server
}

func NewServer(db *repository.DB) *Server {
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

	loadCORS(router)
	loadAuthMiddleware(router, *db)

	return &Server{
		queue:        queue,
		router:       router,
		bundler:      bundler,
		database:     db,
		sessionStore: sessions.NewCookieStore([]byte("siwe-quickstart-secret")),
	}
}

func (s *Server) Start(port string) {
	log.Info().Msgf("connect to http://localhost:%v/ for api/GraphQL playground", port)
	if err := http.ListenAndServe(":"+port, s.router); err != nil {
		log.Fatal().Err(err).Msg("unable to start the server")
	}
}

func (s *Server) Routes() {
	merchantRouter := s.router.PathPrefix("/merchant").Subrouter()

	s.router.Handle("/query", s.walletGraphqlHandler())
	s.router.Handle("/", playground.Handler("api/GraphQL playground", "/query"))

	merchantRouter.Handle("/graphiql", playground.Handler("api/GraphQL playground", "/merchant/query"))
	merchantRouter.Handle("/query", s.merchantGraphqlHandler())

	auth := s.router.PathPrefix("/auth").Subrouter()
	// merchant authentication
	auth.HandleFunc("/nonce", s.GetNonce()).Methods(http.MethodGet)
	auth.HandleFunc("/verify", s.VerifyMerchant())
}

func (s *Server) walletGraphqlHandler() *handler.Server {
	walletSrv := handler.NewDefaultServer(wallet_generated.NewExecutableSchema(wallet_generated.Config{Resolvers: &wallet_graph.Resolver{
		Cache:    repository.NewMCache(),
		Database: s.database,
		// TurnkeyService: tunkeyService,
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
			allowedOrigins := []string{"https://portal.lucidconnect.xyz", "https://checkout.lucidconnect.xyz", "https://lucidconnect.xyz", "https://wallet.lucidconnect.xyz"}
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
			AllowedOrigins: []string{"http://localhost:4002/", "http://localhost:7890/", "http://localhost:3000/"},
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

func loadAuthMiddleware(router *mux.Router, db repository.DB) {
	router.Use(auth.Middleware(db))
}
