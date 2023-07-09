package handlers

import (
	"context"

	"salespot/services/product_service/internal/models"
	"salespot/shared/sctx/component/tracing"
)

type GetProductRepo interface {
	GetProduct(ctx context.Context, id string) (*models.Product, error)
}

type getProductHdl struct {
	repo GetProductRepo
}

func NewGetProductHdl(repo GetProductRepo) *getProductHdl {
	return &getProductHdl{repo: repo}
}

func (h *getProductHdl) Response(ctx context.Context, id string) (*models.Product, error) {
	ctx, span := tracing.StartTrace(ctx, "handler.get-product")
	defer span.End()

	return h.repo.GetProduct(ctx, id)
}
