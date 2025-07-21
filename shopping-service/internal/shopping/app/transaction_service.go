package app

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"shopping-service/internal/shopping/domain"
	"shopping-service/internal/shopping/infra"
)

type TransactionService interface {
	Create(transaction *domain.Transaction) error
	GetAll() ([]domain.Transaction, error)
	GetByID(id string) (*domain.Transaction, error)
	Update(id string, transaction *domain.Transaction) error
	Delete(id string) error
}

type transactionService struct {
	repo infra.TransactionRepository
}

func NewTransactionService(repo infra.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) Create(transaction *domain.Transaction) error {
	// validasi email harus mengandung @
	if !strings.Contains(transaction.Email, "@") {
		return ErrInvalidEmail
	}

	// product_id wajib diisi
	if transaction.ProductID == "" {
		return ErrProductIDMissing
	}

	// quantity harus > 0
	if transaction.Quantity <= 0 {
		return ErrInvalidQuantity
	}

	// hitung total = harga * quantity
	hargaSatuan := 10000.0
	transaction.Total = float64(transaction.Quantity) * hargaSatuan

	// panggil Payment Service
	err := s.callPayment(transaction.Email, transaction.Total, &transaction.PaymentID, &transaction.Status)
	if err != nil {
		return ErrPaymentRequest
	}

	// transaksi gagal jika status bukan success
	if transaction.Status != "success" {
		return ErrPaymentFailed
	}

	// simpan transaksi ke database
	return s.repo.Insert(transaction)
}

func (s *transactionService) callPayment(email string, amount float64, paymentID *string, status *string) error {
	// ambil URL payment dari ENV atau default
	paymentURL := os.Getenv("PAYMENT_URL")
	if paymentURL == "" {
		paymentURL = "http://localhost:8081/payments"
	}

	// buat body request
	payload := map[string]interface{}{
		"email":  email,
		"amount": amount,
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// kirim POST request
	req, err := http.NewRequest("POST", paymentURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// jika status bukan 200, anggap gagal
	if resp.StatusCode != http.StatusOK {
		*status = "failed"
		return nil
	}

	// decode response dari payment service
	var result struct {
		Status    string `json:"status"`
		PaymentID string `json:"payment_id"`
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return err
	}

	*status = result.Status
	*paymentID = result.PaymentID
	return nil
}

// ambil semua transaksi
func (s *transactionService) GetAll() ([]domain.Transaction, error) {
	return s.repo.FindAll()
}

// ambil transaksi by ID
func (s *transactionService) GetByID(id string) (*domain.Transaction, error) {
	return s.repo.FindByID(id)
}

// update transaksi
func (s *transactionService) Update(id string, transaction *domain.Transaction) error {
	return s.repo.Update(id, transaction)
}

// hapus transaksi
func (s *transactionService) Delete(id string) error {
	return s.repo.Delete(id)
}
