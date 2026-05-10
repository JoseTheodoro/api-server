package services

import (
	"apiserver/internal/domain"
	"apiserver/internal/repository"
)

type ProductService interface {
	GetProductFromAnyBusinessRule() *domain.Product
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(r repository.ProductRepository) ProductService {
	return &productService{r}
}

func (p *productService) GetProductFromAnyBusinessRule() *domain.Product {
	return p.repo.Get()
}
