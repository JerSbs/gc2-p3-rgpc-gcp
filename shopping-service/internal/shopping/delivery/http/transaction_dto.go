package http

// struct request transaksi
type CreateTransactionRequest struct {
	Email     string `json:"email"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
