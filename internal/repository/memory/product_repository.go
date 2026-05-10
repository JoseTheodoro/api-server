package memory

import "apiserver/internal/domain"

type ProductRepositoryMemory struct{}

func NewProductRepositoryMemory() *ProductRepositoryMemory {
	return &ProductRepositoryMemory{}
}

func (p *ProductRepositoryMemory) Get() *domain.Product {
	return &domain.Product{
		ID:    1,
		Name:  "Macbook M5",
		Price: 1.800,
	}
}
