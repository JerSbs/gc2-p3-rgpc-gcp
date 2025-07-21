package app

import (
	"context"
	"payment-service/internal/payment/domain"
	"payment-service/internal/payment/infra"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, input domain.Payment) (domain.Payment, error)
	GetPaymentByID(ctx context.Context, id string) (domain.Payment, error)
	DeletePaymentByID(ctx context.Context, id string) (domain.Payment, error)
	GetAllPayments(ctx context.Context) ([]domain.Payment, error)
}

// Implementasi service
type paymentService struct {
	repo infra.PaymentRepository
}

// Inisialisasi service
func NewPaymentService(repo infra.PaymentRepository) PaymentService {
	return &paymentService{repo: repo}
}

func (s *paymentService) CreatePayment(ctx context.Context, input domain.Payment) (domain.Payment, error) {
	// Validasi manual
	if input.Email == "" {
		return domain.Payment{}, ErrEmailEmpty
	}
	if input.Amount <= 0 {
		return domain.Payment{}, ErrAmountInvalid
	}

	// Set default status
	input.Status = "paid"

	// Set ID & Timestamp
	input.ID = primitive.NewObjectID()
	input.CreatedAt = time.Now()

	// Insert ke Mongo
	_, err := s.repo.Insert(ctx, input)
	if err != nil {
		return domain.Payment{}, ErrInsertFailed
	}

	return input, nil
}

// Ambil payment by ID
func (s *paymentService) GetPaymentByID(ctx context.Context, id string) (domain.Payment, error) {
	return s.repo.FindByID(ctx, id)
}

// Hapus payment by ID
func (s *paymentService) DeletePaymentByID(ctx context.Context, id string) (domain.Payment, error) {
	return s.repo.DeleteByID(ctx, id)
}

// Ambil semua payment
func (s *paymentService) GetAllPayments(ctx context.Context) ([]domain.Payment, error) {
	return s.repo.FindAll(ctx)
}
