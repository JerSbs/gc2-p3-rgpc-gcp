package grpc

import (
	"context"
	"payment-service/internal/payment/app"
	"payment-service/internal/payment/delivery/grpc/paymentpb"
	"payment-service/internal/payment/domain"
)

// Handler gRPC
type PaymentHandler struct {
	paymentpb.UnimplementedPaymentServiceServer
	Service app.PaymentService
}

// Tambah payment baru
func (h *PaymentHandler) AddPayment(ctx context.Context, req *paymentpb.AddPaymentRequest) (*paymentpb.Payment, error) {
	input := domain.Payment{
		Email:  req.GetEmail(),
		Amount: req.GetAmount(),
	}

	result, err := h.Service.CreatePayment(ctx, input)
	if err != nil {
		return nil, err
	}

	return &paymentpb.Payment{
		Id:     result.ID.Hex(),
		Email:  result.Email,
		Amount: result.Amount,
		Status: result.Status,
	}, nil
}

// Ambil payment by ID
func (h *PaymentHandler) GetPaymentByID(ctx context.Context, req *paymentpb.GetPaymentByIDRequest) (*paymentpb.Payment, error) {
	result, err := h.Service.GetPaymentByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if result.ID.IsZero() {
		return nil, nil
	}

	return &paymentpb.Payment{
		Id:     result.ID.Hex(),
		Email:  result.Email,
		Amount: result.Amount,
		Status: result.Status,
	}, nil
}

// Hapus payment by ID
func (h *PaymentHandler) DeletePaymentByID(ctx context.Context, req *paymentpb.DeletePaymentByIDRequest) (*paymentpb.Payment, error) {
	result, err := h.Service.DeletePaymentByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &paymentpb.Payment{
		Id:     result.ID.Hex(),
		Email:  result.Email,
		Amount: result.Amount,
		Status: result.Status,
	}, nil
}

// Ambil semua payment
func (h *PaymentHandler) GetAllPayments(ctx context.Context, _ *paymentpb.GetAllPaymentsRequest) (*paymentpb.GetAllPaymentsResponse, error) {
	data, err := h.Service.GetAllPayments(ctx)
	if err != nil {
		return nil, err
	}

	var payments []*paymentpb.Payment
	for _, p := range data {
		payments = append(payments, &paymentpb.Payment{
			Id:     p.ID.Hex(),
			Email:  p.Email,
			Amount: p.Amount,
			Status: p.Status,
		})
	}

	return &paymentpb.GetAllPaymentsResponse{Payments: payments}, nil
}
