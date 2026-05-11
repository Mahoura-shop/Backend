package shipping

import "github.com/Mahoura-shop/Backend/internal/domain/entity"

const (
	DefaultShippingCost = 50000
)

func CalculateShipping(province *entity.Province) uint {
	if province == nil || province.ShippingCost == 0 {
		return DefaultShippingCost
	}
	return province.ShippingCost
}
