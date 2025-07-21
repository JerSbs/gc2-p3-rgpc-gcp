package http

import "github.com/labstack/echo/v4"

// setup route transaksi
func TransactionRoute(e *echo.Echo, handler *TransactionHandler) {
	route := e.Group("/transactions")

	route.POST("", handler.CreateTransaction)       // tambah transaksi
	route.GET("", handler.GetAllTransactions)       // ambil semua
	route.GET("/:id", handler.GetTransactionByID)   // ambil by id
	route.PUT("/:id", handler.UpdateTransaction)    // update by id
	route.DELETE("/:id", handler.DeleteTransaction) // hapus by id
}
