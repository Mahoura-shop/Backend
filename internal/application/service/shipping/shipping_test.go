package shipping_test

import (
	"testing"

	"github.com/Mahoura-shop/Backend/internal/application/service/shipping"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
)

func TestCalculateShipping_NilProvince(t *testing.T) {
	got := shipping.CalculateShipping(nil)
	if got != shipping.DefaultShippingCost {
		t.Errorf("got %d, want %d", got, shipping.DefaultShippingCost)
	}
}

func TestCalculateShipping_ZeroShippingCost(t *testing.T) {
	p := &entity.Province{ShippingCost: 0}
	got := shipping.CalculateShipping(p)
	if got != shipping.DefaultShippingCost {
		t.Errorf("got %d, want %d", got, shipping.DefaultShippingCost)
	}
}

func TestCalculateShipping_CustomCost(t *testing.T) {
	p := &entity.Province{ShippingCost: 30000}
	got := shipping.CalculateShipping(p)
	if got != 30000 {
		t.Errorf("got %d, want 30000", got)
	}
}

func TestCalculateShipping_DefaultConstantValue(t *testing.T) {
	if shipping.DefaultShippingCost == 0 {
		t.Error("DefaultShippingCost should not be zero")
	}
}
