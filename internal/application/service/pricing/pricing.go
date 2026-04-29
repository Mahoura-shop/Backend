package pricing

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
)

func ResolvePrice(userType enum.UserType, product entity.Product) uint {
	switch userType {
	case enum.UserTypeFellow:
		return product.Step1Price
	case enum.UserTypeShopkeeperCash:
		return product.Step2Price
	case enum.UserTypeShopkeeperCheque:
		return product.Step3Price
	case enum.UserTypeCustomer:
		return product.Step4Price
	default:
		return product.ConsumerPrice
	}
}
