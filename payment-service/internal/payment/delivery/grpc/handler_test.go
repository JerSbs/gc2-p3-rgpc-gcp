package grpc

import (
	"context"
	"payment-service/internal/payment/delivery/grpc/paymentpb"
	"payment-service/internal/payment/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock service
type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) CreatePayment(ctx context.Context, input domain.Payment) (domain.Payment, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(domain.Payment), args.Error(1)
}
func (m *MockPaymentService) GetPaymentByID(ctx context.Context, id string) (domain.Payment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Payment), args.Error(1)
}
func (m *MockPaymentService) DeletePaymentByID(ctx context.Context, id string) (domain.Payment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Payment), args.Error(1)
}
func (m *MockPaymentService) GetAllPayments(ctx context.Context) ([]domain.Payment, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Payment), args.Error(1)
}

//  TESTS

func TestAddPayment_Success(t *testing.T) {
	mockSvc := new(MockPaymentService)
	handler := &PaymentHandler{Service: mockSvc}

	input := &paymentpb.AddPaymentRequest{Email: "user@example.com", Amount: 100}
	fakeID := primitive.NewObjectID()
	expected := domain.Payment{
		ID:        fakeID,
		Email:     "user@example.com",
		Amount:    100,
		Status:    "paid",
		CreatedAt: time.Now(),
	}
	mockSvc.On("CreatePayment", mock.Anything, mock.AnythingOfType("domain.Payment")).Return(expected, nil)

	resp, err := handler.AddPayment(context.TODO(), input)

	assert.NoError(t, err)
	assert.Equal(t, expected.Email, resp.Email)
	assert.Equal(t, expected.Amount, resp.Amount)
	assert.Equal(t, expected.Status, resp.Status)
	assert.Equal(t, expected.ID.Hex(), resp.Id)
	mockSvc.AssertExpectations(t)
}

func TestGetPaymentByID_Success(t *testing.T) {
	mockSvc := new(MockPaymentService)
	handler := &PaymentHandler{Service: mockSvc}

	fakeID := primitive.NewObjectID()
	expected := domain.Payment{ID: fakeID, Email: "x@y.com", Amount: 55.5, Status: "paid"}
	mockSvc.On("GetPaymentByID", mock.Anything, "abc123").Return(expected, nil)

	resp, err := handler.GetPaymentByID(context.TODO(), &paymentpb.GetPaymentByIDRequest{Id: "abc123"})

	assert.NoError(t, err)
	assert.Equal(t, expected.ID.Hex(), resp.Id)
	assert.Equal(t, expected.Email, resp.Email)
	assert.Equal(t, expected.Amount, resp.Amount)
	assert.Equal(t, expected.Status, resp.Status)
}

func TestDeletePaymentByID_Success(t *testing.T) {
	mockSvc := new(MockPaymentService)
	handler := &PaymentHandler{Service: mockSvc}

	fakeID := primitive.NewObjectID()
	expected := domain.Payment{ID: fakeID, Email: "del@x.com", Amount: 99, Status: "deleted"}
	mockSvc.On("DeletePaymentByID", mock.Anything, "del123").Return(expected, nil)

	resp, err := handler.DeletePaymentByID(context.TODO(), &paymentpb.DeletePaymentByIDRequest{Id: "del123"})

	assert.NoError(t, err)
	assert.Equal(t, expected.ID.Hex(), resp.Id)
	assert.Equal(t, expected.Email, resp.Email)
	assert.Equal(t, expected.Amount, resp.Amount)
	assert.Equal(t, expected.Status, resp.Status)
}

func TestGetAllPayments_Success(t *testing.T) {
	mockSvc := new(MockPaymentService)
	handler := &PaymentHandler{Service: mockSvc}

	fakeID := primitive.NewObjectID()
	expected := []domain.Payment{
		{ID: fakeID, Email: "a@a.com", Amount: 1, Status: "paid"},
	}
	mockSvc.On("GetAllPayments", mock.Anything).Return(expected, nil)

	resp, err := handler.GetAllPayments(context.TODO(), &paymentpb.GetAllPaymentsRequest{})

	assert.NoError(t, err)
	assert.Len(t, resp.Payments, 1)
	assert.Equal(t, "a@a.com", resp.Payments[0].Email)
	assert.Equal(t, float64(1), resp.Payments[0].Amount)
	assert.Equal(t, "paid", resp.Payments[0].Status)
}
