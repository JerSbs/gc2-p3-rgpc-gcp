package http

import "github.com/labstack/echo/v4"

// setup route
func ProductRoute(e *echo.Echo, handler *ProductHandler) {
	route := e.Group("/products")

	route.POST("", handler.CreateProduct)       // tambah product
	route.GET("", handler.GetAllProducts)       // ambil semua
	route.GET("/:id", handler.GetProductByID)   // ambil by id
	route.PUT("/:id", handler.UpdateProduct)    // update by id
	route.DELETE("/:id", handler.DeleteProduct) // hapus by id
}
