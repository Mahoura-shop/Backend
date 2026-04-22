package orderdto

type CreateOrderRequest struct {
	UserID uint
}

type CreateOrderItemRequest struct {
	ProductID     uint
	Count         uint
	PriceSnapshot uint
}