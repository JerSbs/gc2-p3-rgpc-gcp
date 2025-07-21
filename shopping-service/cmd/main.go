package main

import (
	"log"
	"shopping-service/config"

	"shopping-service/internal/shopping/app"
	"shopping-service/internal/shopping/delivery/http"
	"shopping-service/internal/shopping/infra"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "shopping-service/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Shopping Service API
// @version 1.0
// @description REST API untuk produk dan transaksi
// @host localhost:8080
// @BasePath /
func main() {
	// load env & koneksi MongoDB
	config.LoadEnv()
	config.ConnectDB()

	// init Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// init product
	productRepo := infra.NewProductRepo()
	productService := app.NewProductService(productRepo)
	productHandler := http.NewProductHandler(productService)
	http.ProductRoute(e, productHandler)

	// init transaction
	transactionRepo := infra.NewTransactionRepo()
	transactionService := app.NewTransactionService(transactionRepo)
	transactionHandler := http.NewTransactionHandler(&transactionService)
	http.TransactionRoute(e, transactionHandler)

	// Jalankan cron job transaksi
	app.StartTransactionCron(transactionRepo)

	// start server
	port := config.GetEnv("PORT")
	log.Println("running at port:", port)
	e.Logger.Fatal(e.Start(":" + port))
}
