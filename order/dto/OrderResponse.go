package dto

type OrderResponse struct {
	OrderCode  string
	Products   []ProductDetail
	TotalPrice int64
}

type ProductDetail struct {
	ID       int64
	Name     string
	Quantity int32
	Price    int64
}
