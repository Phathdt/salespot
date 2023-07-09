package handlers

import (
	"context"

	"salespot/services/product_service/internal/models"
)

type listProductRepo interface {
	ListProduct(ctx context.Context) ([]models.Product, error)
}

type listProductHdl struct {
	repo listProductRepo
}

func NewListProductHdl(repo listProductRepo) *listProductHdl {
	return &listProductHdl{repo: repo}
}

func (h *listProductHdl) Response(ctx context.Context) ([]models.Product, error) {
	return h.repo.ListProduct(ctx)
}
