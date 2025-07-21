package app

import "errors"

// validasi input
var (
	ErrNameEmpty        = errors.New("product name is required")
	ErrPriceInvalid     = errors.New("price must be > 0")
	ErrStockInvalid     = errors.New("stock must be >= 0")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrInvalidQuantity  = errors.New("quantity must be > 0")
	ErrProductIDMissing = errors.New("product_id is required")
)

// insert
var (
	ErrProductInsert     = errors.New("failed to insert product")
	ErrTransactionInsert = errors.New("failed to insert transaction")
)

// read
var (
	ErrInvalidProductID     = errors.New("invalid product ID")
	ErrInvalidTransactionID = errors.New("invalid transaction ID")
	ErrProductNotFound      = errors.New("product not found")
	ErrTransactionNotFound  = errors.New("transaction not found")
	ErrFailedDecode         = errors.New("failed to decode data")
)

// update
var (
	ErrProductUpdate     = errors.New("failed to update product")
	ErrTransactionUpdate = errors.New("failed to update transaction")
)

// delete
var (
	ErrProductDelete     = errors.New("failed to delete product")
	ErrTransactionDelete = errors.New("failed to delete transaction")
)

// transaksi / payment
var (
	ErrInvalidTransaction = errors.New("invalid transaction input")
	ErrPaymentFailed      = errors.New("payment failed")
	ErrPaymentRequest     = errors.New("payment service not reachable")
)
