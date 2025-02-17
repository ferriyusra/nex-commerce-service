package app

import (
	"context"
	"log"
	"nex-commerce-service/config"
	"nex-commerce-service/internal/adapter/handler"
	"nex-commerce-service/internal/adapter/repository"
	"nex-commerce-service/internal/core/service"
	"nex-commerce-service/lib/auth"
	"nex-commerce-service/lib/middleware"
	"nex-commerce-service/lib/pagination"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func RunServer() {
	cfg := config.NewConfig()
	db, err := cfg.ConnectionPostgres()
	if err != nil {
		log.Fatal("error connection database %v", err)
		return
	}

	err = os.MkdirAll("./temp/content", 0755)
	if err != nil {
		log.Fatal("error creating temp directory %v", err)
		return
	}

	// auth and middleware
	jwt := auth.NewJwt(cfg)
	middlewareAuth := middleware.NewMiddleware(cfg)

	// pagination
	_ = pagination.NewPagination()

	// repository
	authRepo := repository.NewAuthRepository(db.DB)
	accountRepo := repository.NewAccountRepository(db.DB)
	productRepo := repository.NewProductRepository(db.DB)
	financialRepo := repository.NewFinancialRepository(db.DB)
	cartRepo := repository.NewCartRepository(db.DB)
	transactionRepo := repository.NewTransactionRepository(db.DB)

	// service
	authService := service.NewAuthService(authRepo, accountRepo, cfg, jwt)
	productService := service.NewProductService(productRepo)
	financialService := service.NewFinancialService(financialRepo)
	cartService := service.NewCartService(cartRepo)
	transactionService := service.NewTransactionService(transactionRepo)

	// handler
	authHandler := handler.NewAuthHandler(authService)
	productHandler := handler.NewProductHandler(productService)
	financialHandler := handler.NewFinancialHandler(financialService)
	cartHandler := handler.NewCartHandler(cartService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// intitalization server
	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip} ${status} - ${latency} ${method} ${path}\n",
	}))

	api := app.Group("/api")
	api.Post("/auth/login", authHandler.Login)

	// user as customer
	customerApp := api.Group("/auth/customers")
	customerApp.Post("/register", authHandler.RegisterCustomer)

	// user as seller
	sellerApp := api.Group("/auth/sellers")
	sellerApp.Post("/register", authHandler.RegisterSeller)

	productApp := api.Group("/products")
	productApp.Get("/", productHandler.FindAll)

	financialApp := api.Group("/financial")
	financialApp.Use(middlewareAuth.CheckToken())
	financialApp.Get("/balance", middleware.ACLMiddleware([]string{"customer", "seller"}), financialHandler.GetBalance)
	financialApp.Post("/deposit", middleware.ACLMiddleware([]string{"customer"}), financialHandler.Deposit)
	financialApp.Post("/withdraw", middleware.ACLMiddleware([]string{"seller"}), financialHandler.Withdraw)

	cartApp := api.Group("/carts")
	cartApp.Use(middlewareAuth.CheckToken())
	cartApp.Get("/", middleware.ACLMiddleware([]string{"customer"}), cartHandler.GetCartByUserID)
	cartApp.Post("/add", middleware.ACLMiddleware([]string{"customer"}), cartHandler.AddToCart)

	transactionApp := api.Group("/transactions")
	transactionApp.Use(middlewareAuth.CheckToken())
	transactionApp.Post("/checkout", middleware.ACLMiddleware([]string{"customer"}), transactionHandler.Checkout)

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err := app.Listen(":" + cfg.App.AppPort)
		if err != nil {
			log.Fatal("Error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit

	log.Println("server shutdown of 5 seconds")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)
}
