package models

import (
	"salespot/shared/sctx/core"
)

type Product struct {
	core.MgoModel `json:",inline" bson:",inline"`
	Name          string `json:"name" bson:"name"`
	Price         int    `json:"price" bson:"price"`
}

func (p Product) Collection() string {
	return "products"
}
