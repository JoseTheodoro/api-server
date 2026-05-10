package repository

import "apiserver/internal/domain"

type ProductRepository interface {
	Get() *domain.Product
}
