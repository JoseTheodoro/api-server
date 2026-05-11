package handlers

import (
	"log/slog"
	"net/http"

	"apiserver/internal/services"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(s services.ProductService) *ProductHandler {
	return &ProductHandler{s}
}

func (p *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	slog.Info("request product recieved")
	product := p.productService.GetProductFromAnyBusinessRule()
	slog.Info("product got", "product", product)

	if err := writeJSON(w, http.StatusOK, product); err != nil {
		slog.Error("unable to encode to JSON", "err", err)
	}
}
