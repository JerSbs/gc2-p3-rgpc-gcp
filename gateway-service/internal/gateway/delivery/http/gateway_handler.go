package http

import (
	"context"
	"net/http"

	"gateway-service/config"
	"gateway-service/internal/gateway/infra"
	"gateway-service/proto"

	"github.com/labstack/echo/v4"
)

// GatewayHandler menangani permintaan REST/gRPC ke semua service
type GatewayHandler struct{}

func NewGatewayHandler() *GatewayHandler {
	return &GatewayHandler{}
}

//  PROXY FORWARD

func (h *GatewayHandler) AuthProxy(c echo.Context) error {
	targetURL := config.AuthServiceURL + c.Path()
	return infra.ForwardRequest(c, targetURL)
}

func (h *GatewayHandler) ShoppingProxy(c echo.Context) error {
	targetURL := config.ShoppingServiceURL + c.Path()
	return infra.ForwardRequest(c, targetURL)
}

func (h *GatewayHandler) PaymentProxy(c echo.Context) error {
	targetURL := config.PaymentServiceURL + c.Path()
	return infra.ForwardRequest(c, targetURL)
}

//  AUTH HANDLER

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *GatewayHandler) LoginHandler(c echo.Context) error {
	var input LoginInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid input"})
	}

	client := infra.GetAuthClient()
	ctx := context.Background()

	req := &proto.LoginRequest{
		Email:    input.Email,
		Password: input.Password,
	}

	res, err := client.LoginUser(ctx, req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Login failed", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": res.Token})
}

func (h *GatewayHandler) RegisterHandler(c echo.Context) error {
	var input RegisterInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid input"})
	}

	client := infra.GetAuthClient()
	ctx := context.Background()

	req := &proto.RegisterRequest{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	res, err := client.RegisterUser(ctx, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Register failed", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id":    res.Id,
		"name":  res.Name,
		"email": res.Email,
	})
}

//  PAYMENT HANDLER (gRPC)

func (h *GatewayHandler) CreatePaymentHandler(c echo.Context) error {
	var input struct {
		Email  string `json:"email"`
		Amount int64  `json:"amount"`
		Status string `json:"status"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid input format"})
	}

	if input.Email == "" || input.Amount <= 0 || input.Status == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "email, amount, and status are required"})
	}

	client := infra.GetPaymentClient()
	ctx := context.Background()

	req := &proto.AddPaymentRequest{
		Email:  input.Email,
		Amount: float64(input.Amount),
	}

	res, err := client.AddPayment(ctx, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to create payment", "error": err.Error()})
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *GatewayHandler) GetAllPaymentsHandler(c echo.Context) error {
	client := infra.GetPaymentClient()
	ctx := context.Background()

	res, err := client.GetAllPayments(ctx, &proto.GetAllPaymentsRequest{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to fetch payments", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, res.Payments)
}

func (h *GatewayHandler) GetPaymentByIDHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Missing payment ID"})
	}

	client := infra.GetPaymentClient()
	ctx := context.Background()

	res, err := client.GetPaymentByID(ctx, &proto.GetPaymentByIDRequest{Id: id})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to fetch payment", "error": err.Error()})
	}
	if res == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Payment not found"})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *GatewayHandler) DeletePaymentByIDHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Missing payment ID"})
	}

	client := infra.GetPaymentClient()
	ctx := context.Background()

	res, err := client.DeletePaymentByID(ctx, &proto.DeletePaymentByIDRequest{Id: id})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to delete payment", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

//  SHOPPING - PRODUCT (REST PROXY)

func (h *GatewayHandler) GetAllProducts(c echo.Context) error {
	return infra.ForwardRequest(c, config.ShoppingServiceURL+"/products")
}

func (h *GatewayHandler) CreateProduct(c echo.Context) error {
	return infra.ForwardRequest(c, config.ShoppingServiceURL+"/products")
}

func (h *GatewayHandler) GetProductByID(c echo.Context) error {
	id := c.Param("id")
	return infra.ForwardRequest(c, config.ShoppingServiceURL+"/products/"+id)
}

func (h *GatewayHandler) UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	return infra.ForwardRequest(c, config.ShoppingServiceURL+"/products/"+id)
}

func (h *GatewayHandler) DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	return infra.ForwardRequest(c, config.ShoppingServiceURL+"/products/"+id)
}

//  SHOPPING - TRANSACTION (REST PROXY)

func (h *GatewayHandler) GetAllTransactions(c echo.Context) error {
	return infra.ForwardRequest(c, config.ShoppingServiceURL+"/transactions")
}

func (h *GatewayHandler) CreateTransaction(c echo.Context) error {
	return infra.ForwardRequest(c, config.ShoppingServiceURL+"/transactions")
}

func (h *GatewayHandler) GetTransactionByID(c echo.Context) error {
	id := c.Param("id")
	return infra.ForwardRequest(c, config.ShoppingServiceURL+"/transactions/"+id)
}

func (h *GatewayHandler) UpdateTransaction(c echo.Context) error {
	id := c.Param("id")
	return infra.ForwardRequest(c, config.ShoppingServiceURL+"/transactions/"+id)
}

func (h *GatewayHandler) DeleteTransaction(c echo.Context) error {
	id := c.Param("id")
	return infra.ForwardRequest(c, config.ShoppingServiceURL+"/transactions/"+id)
}
