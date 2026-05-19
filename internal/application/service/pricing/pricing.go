package pricing

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
)

func firstNonZero(vals ...uint) uint {
	for _, v := range vals {
		if v != 0 {
			return v
		}
	}
	return 0
}

func ResolvePrice(userType enum.UserType, product entity.Product) uint {
	switch userType {
	case enum.UserTypeFellow, enum.UserTypeAdmin:
		return firstNonZero(product.Step1Price, product.Step2Price, product.Step3Price, product.Step4Price, product.ConsumerPrice, product.IRRPrice)
	case enum.UserTypeShopkeeperCash, enum.UserTypeShopkeeperCheque:
		return firstNonZero(product.Step2Price, product.Step3Price, product.Step4Price, product.ConsumerPrice, product.IRRPrice)
	case enum.UserTypeCustomer:
		return firstNonZero(product.Step4Price, product.ConsumerPrice, product.IRRPrice)
	default:
		return firstNonZero(product.ConsumerPrice, product.IRRPrice)
	}
}
