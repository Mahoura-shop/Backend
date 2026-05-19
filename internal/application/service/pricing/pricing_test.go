package pricing_test

import (
	"testing"

	"github.com/Mahoura-shop/Backend/internal/application/service/pricing"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
)

var testProduct = entity.Product{
	Step1Price:    1_000_000,
	Step2Price:    2_000_000,
	Step3Price:    3_000_000,
	Step4Price:    4_000_000,
	ConsumerPrice: 5_000_000,
}

func TestResolvePrice(t *testing.T) {
	cases := []struct {
		userType enum.UserType
		expected uint
	}{
		{enum.UserTypeFellow, 1_000_000},
		{enum.UserTypeShopkeeperCash, 2_000_000},
		{enum.UserTypeShopkeeperCheque, 2_000_000},
		{enum.UserTypeCustomer, 4_000_000},
		{enum.UserTypeAdmin, 1_000_000},
	}

	for _, c := range cases {
		got := pricing.ResolvePrice(c.userType, testProduct)
		if got != c.expected {
			t.Errorf("ResolvePrice(%v) = %d, want %d", c.userType, got, c.expected)
		}
	}
}

func TestResolvePrice_ZeroPrices(t *testing.T) {
	emptyProduct := entity.Product{}
	userTypes := []enum.UserType{
		enum.UserTypeFellow,
		enum.UserTypeShopkeeperCash,
		enum.UserTypeShopkeeperCheque,
		enum.UserTypeCustomer,
	}
	for _, ut := range userTypes {
		got := pricing.ResolvePrice(ut, emptyProduct)
		if got != 0 {
			t.Errorf("ResolvePrice on empty product with %v = %d, want 0", ut, got)
		}
	}
}

func TestResolvePrice_IRRPriceFallback(t *testing.T) {
	irrOnlyProduct := entity.Product{IRRPrice: 500_000}
	userTypes := []enum.UserType{
		enum.UserTypeFellow,
		enum.UserTypeShopkeeperCash,
		enum.UserTypeShopkeeperCheque,
		enum.UserTypeCustomer,
		enum.UserTypeAdmin,
	}
	for _, ut := range userTypes {
		got := pricing.ResolvePrice(ut, irrOnlyProduct)
		if got != 500_000 {
			t.Errorf("ResolvePrice IRR fallback with %v = %d, want 500000", ut, got)
		}
	}
}

func TestResolvePrice_Ordering(t *testing.T) {
	p := testProduct
	fellow := pricing.ResolvePrice(enum.UserTypeFellow, p)
	cash := pricing.ResolvePrice(enum.UserTypeShopkeeperCash, p)
	cheque := pricing.ResolvePrice(enum.UserTypeShopkeeperCheque, p)
	customer := pricing.ResolvePrice(enum.UserTypeCustomer, p)

	if !(fellow <= cash && cash <= cheque && cheque <= customer) {
		t.Errorf("price ordering violated: fellow=%d cash=%d cheque=%d customer=%d",
			fellow, cash, cheque, customer)
	}
}
