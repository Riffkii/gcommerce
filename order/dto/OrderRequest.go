package dto

type OrderRequest struct {
	CustomerId int64     `json:"customerId"`
	Products   []Product `json:"products"`
}

type Product struct {
	ProductId int64 `json:"productId"`
	Quantity  int32 `json:"quantity"`
}
