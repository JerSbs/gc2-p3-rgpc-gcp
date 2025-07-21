package app

import (
	"shopping-service/internal/shopping/domain"
	"shopping-service/internal/shopping/infra"
)

// service product
type ProductService struct {
	Repo *infra.ProductRepo
}

// init service
func NewProductService(repo *infra.ProductRepo) *ProductService {
	return &ProductService{Repo: repo}
}

// logika simpan product
func (s *ProductService) CreateProduct(product *domain.Product) error {
	// validasi manual
	if product.Name == "" {
		return ErrNameEmpty
	}
	if product.Price <= 0 {
		return ErrPriceInvalid
	}
	if product.Stock < 0 {
		return ErrStockInvalid
	}

	// simpan ke repo
	err := s.Repo.CreateProduct(product)
	if err != nil {
		return ErrProductInsert
	}

	return nil
}

// ambil semua product
func (s *ProductService) GetAllProducts() ([]domain.Product, error) {
	products, err := s.Repo.GetAllProducts()
	if err != nil {
		return nil, ErrFailedDecode
	}
	return products, nil
}

// ambil product by ID
func (s *ProductService) GetProductByID(id string) (*domain.Product, error) {
	product, err := s.Repo.GetProductByID(id)
	if err != nil {
		// jika format ObjectID salah
		if err.Error() == "invalid product ID" {
			return nil, ErrInvalidProductID
		}
		// jika data tidak ditemukan
		return nil, ErrProductNotFound
	}
	return product, nil
}

// update product by ID
func (s *ProductService) UpdateProduct(id string, product *domain.Product) error {
	// validasi
	if product.Name == "" {
		return ErrNameEmpty
	}
	if product.Price <= 0 {
		return ErrPriceInvalid
	}
	if product.Stock < 0 {
		return ErrStockInvalid
	}

	err := s.Repo.UpdateProduct(id, product)
	if err != nil {
		if err.Error() == "invalid product ID" {
			return ErrInvalidProductID
		}
		return ErrProductUpdate
	}

	return nil
}

// delete product by ID
func (s *ProductService) DeleteProduct(id string) error {
	err := s.Repo.DeleteProduct(id)
	if err != nil {
		if err.Error() == "invalid product ID" {
			return ErrInvalidProductID
		}
		return ErrProductDelete
	}
	return nil
}
