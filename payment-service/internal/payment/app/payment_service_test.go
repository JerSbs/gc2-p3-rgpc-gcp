package app

import (
	"context"
	"errors"
	"payment-service/internal/payment/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mock repository
type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Insert(ctx context.Context, payment domain.Payment) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, payment)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockPaymentRepository) FindByID(ctx context.Context, id string) (domain.Payment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Payment), args.Error(1)
}

func (m *MockPaymentRepository) DeleteByID(ctx context.Context, id string) (domain.Payment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Payment), args.Error(1)
}

func (m *MockPaymentRepository) FindAll(ctx context.Context) ([]domain.Payment, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Payment), args.Error(1)
}

// ===================== TESTS ========================

func TestCreatePayment_Success(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewPaymentService(mockRepo)

	input := domain.Payment{
		Email:  "user@example.com",
		Amount: 100.0,
	}

	mockRepo.On("Insert", mock.Anything, mock.MatchedBy(func(p domain.Payment) bool {
		return p.Email == input.Email && p.Amount == input.Amount && p.Status == "paid"
	})).Return(&mongo.InsertOneResult{}, nil)

	result, err := service.CreatePayment(context.TODO(), input)

	assert.NoError(t, err)
	assert.Equal(t, input.Email, result.Email)
	assert.Equal(t, input.Amount, result.Amount)
	assert.Equal(t, "paid", result.Status)
	assert.NotEmpty(t, result.ID)
	assert.WithinDuration(t, time.Now(), result.CreatedAt, 2*time.Second)
	mockRepo.AssertExpectations(t)
}

func TestCreatePayment_EmptyEmail(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewPaymentService(mockRepo)

	input := domain.Payment{Email: "", Amount: 100.0}
	result, err := service.CreatePayment(context.TODO(), input)

	assert.ErrorIs(t, err, ErrEmailEmpty)
	assert.Equal(t, domain.Payment{}, result)
}

func TestCreatePayment_InvalidAmount(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewPaymentService(mockRepo)

	input := domain.Payment{Email: "x@y.com", Amount: 0}
	result, err := service.CreatePayment(context.TODO(), input)

	assert.ErrorIs(t, err, ErrAmountInvalid)
	assert.Equal(t, domain.Payment{}, result)
}

func TestCreatePayment_InsertFailed(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewPaymentService(mockRepo)

	input := domain.Payment{Email: "user@example.com", Amount: 100.0}
	mockRepo.On("Insert", mock.Anything, mock.Anything).Return(nil, errors.New("mongo error"))

	result, err := service.CreatePayment(context.TODO(), input)

	assert.ErrorIs(t, err, ErrInsertFailed)
	assert.Equal(t, domain.Payment{}, result)
}

func TestGetPaymentByID_Success(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewPaymentService(mockRepo)

	expected := domain.Payment{Email: "x@y.com", Amount: 42}
	mockRepo.On("FindByID", mock.Anything, "abc123").Return(expected, nil)

	result, err := service.GetPaymentByID(context.TODO(), "abc123")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestDeletePaymentByID_Success(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewPaymentService(mockRepo)

	expected := domain.Payment{Email: "del@x.com", Amount: 500}
	mockRepo.On("DeleteByID", mock.Anything, "del123").Return(expected, nil)

	result, err := service.DeletePaymentByID(context.TODO(), "del123")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAllPayments_Success(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewPaymentService(mockRepo)

	expected := []domain.Payment{
		{Email: "a@a.com", Amount: 1},
		{Email: "b@b.com", Amount: 2},
	}
	mockRepo.On("FindAll", mock.Anything).Return(expected, nil)

	result, err := service.GetAllPayments(context.TODO())

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}
