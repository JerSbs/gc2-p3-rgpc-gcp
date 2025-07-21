package http

import (
	"net/http"
	"shopping-service/internal/shopping/app"
	"shopping-service/internal/shopping/domain"

	"github.com/labstack/echo/v4"
)

// struct handler
type ProductHandler struct {
	Service *app.ProductService
}

// init handler
func NewProductHandler(service *app.ProductService) *ProductHandler {
	return &ProductHandler{Service: service}
}

// CreateProduct godoc
// @Summary Tambah product baru
// @Description Menambahkan product ke database
// @Tags Products
// @Accept json
// @Produce json
// @Param request body CreateProductRequest true "Product data"
// @Success 201 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var req CreateProductRequest

	if err := c.Bind(&req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "invalid request format")
	}

	product := &domain.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}

	err := h.Service.CreateProduct(product)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"message": "product created",
		"data": map[string]any{
			"id":         product.ID.Hex(),
			"name":       product.Name,
			"price":      product.Price,
			"stock":      product.Stock,
			"created_at": product.CreatedAt,
		},
	})
}

// GetAllProducts godoc
// @Summary Ambil semua product
// @Description Menampilkan seluruh daftar product
// @Tags Products
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c echo.Context) error {
	products, err := h.Service.GetAllProducts()
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"data": products,
	})
}

// GetProductByID godoc
// @Summary Ambil product berdasarkan ID
// @Description Menampilkan data product spesifik berdasarkan ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c echo.Context) error {
	id := c.Param("id")

	product, err := h.Service.GetProductByID(id)
	if err != nil {
		return ErrorResponse(c, http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"data": map[string]any{
			"id":         product.ID.Hex(),
			"name":       product.Name,
			"price":      product.Price,
			"stock":      product.Stock,
			"created_at": product.CreatedAt,
		},
	})
}

// UpdateProduct godoc
// @Summary Update product berdasarkan ID
// @Description Update data product spesifik
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body CreateProductRequest true "Updated product data"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	var req CreateProductRequest

	if err := c.Bind(&req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "invalid request format")
	}

	product := &domain.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}

	err := h.Service.UpdateProduct(id, product)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "product updated",
	})
}

// DeleteProduct godoc
// @Summary Hapus product berdasarkan ID
// @Description Menghapus data product dari database
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id := c.Param("id")

	err := h.Service.DeleteProduct(id)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "product deleted",
	})
}
