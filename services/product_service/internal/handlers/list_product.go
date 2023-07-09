package handlers

import (
	"context"

	"salespot/services/product_service/internal/models"
	"salespot/shared/sctx/component/tracing"
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
	ctx, span := tracing.StartTrace(ctx, "handler.list-product")
	defer span.End()

	return h.repo.ListProduct(ctx)
}
