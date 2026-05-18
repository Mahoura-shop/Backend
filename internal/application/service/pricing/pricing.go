package pricing

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
)

func ResolvePrice(userType enum.UserType, product entity.Product) uint {
	switch userType {
	case enum.UserTypeFellow:
		return product.Step1Price
	case enum.UserTypeShopkeeperCash, enum.UserTypeShopkeeperCheque:
		return product.Step2Price
	case enum.UserTypeCustomer:
		return product.Step4Price
	case enum.UserTypeAdmin:
		return product.Step1Price
	default:
		return product.ConsumerPrice
	}
}
