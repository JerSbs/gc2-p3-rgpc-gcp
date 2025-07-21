package http

import (
	"net/http"

	"shopping-service/internal/shopping/app"
	"shopping-service/internal/shopping/domain"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// struct handler transaksi
type TransactionHandler struct {
	Service *app.TransactionService
}

// init handler transaksi
func NewTransactionHandler(service *app.TransactionService) *TransactionHandler {
	return &TransactionHandler{Service: service}
}

// CreateTransaction godoc
// @Summary Tambah transaksi
// @Description Menambahkan transaksi dan memanggil Payment Service
// @Tags Transactions
// @Accept json
// @Produce json
// @Param request body CreateTransactionRequest true "Transaksi data"
// @Success 201 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	var req CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	if req.Email == "" || req.ProductID == "" || req.Quantity <= 0 {
		return ErrorResponse(c, http.StatusBadRequest, "email, product_id, and quantity are required")
	}

	tx := &domain.Transaction{
		Email:     req.Email,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := (*h.Service).Create(tx); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"message": "transaction created",
		"data":    tx,
	})
}

// GetAllTransactions godoc
// @Summary Ambil semua transaksi
// @Description Menampilkan seluruh transaksi
// @Tags Transactions
// @Produce json
// @Success 200 {object} []domain.Transaction
// @Failure 500 {object} map[string]any
// @Router /transactions [get]
func (h *TransactionHandler) GetAllTransactions(c echo.Context) error {
	result, err := (*h.Service).GetAll()
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

// GetTransactionByID godoc
// @Summary Ambil transaksi berdasarkan ID
// @Description Menampilkan transaksi spesifik berdasarkan ID
// @Tags Transactions
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} domain.Transaction
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransactionByID(c echo.Context) error {
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "invalid transaction id")
	}

	result, err := (*h.Service).GetByID(id)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	if result == nil {
		return ErrorResponse(c, http.StatusNotFound, "transaction not found")
	}

	return c.JSON(http.StatusOK, result)
}

// UpdateTransaction godoc
// @Summary Update transaksi berdasarkan ID
// @Description Update data transaksi tertentu
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Param request body CreateTransactionRequest true "Updated transaksi data"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /transactions/{id} [put]
func (h *TransactionHandler) UpdateTransaction(c echo.Context) error {
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "invalid transaction id")
	}

	var req CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}
	if req.Email == "" || req.ProductID == "" || req.Quantity <= 0 {
		return ErrorResponse(c, http.StatusBadRequest, "email, product_id, and quantity are required")
	}

	tx := &domain.Transaction{
		Email:     req.Email,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := (*h.Service).Update(id, tx); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "transaction updated",
	})
}

// DeleteTransaction godoc
// @Summary Hapus transaksi berdasarkan ID
// @Description Menghapus transaksi dari database
// @Tags Transactions
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /transactions/{id} [delete]
func (h *TransactionHandler) DeleteTransaction(c echo.Context) error {
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "invalid transaction id")
	}

	if err := (*h.Service).Delete(id); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "transaction deleted",
	})
}
