package models

type Product struct {
	ID    string `json:"id" bson:"_id"`
	Name  string `json:"name" bson:"name"`
	Price int    `json:"price" bson:"price"`
}

func (p Product) Collection() string {
	return "products"
}
