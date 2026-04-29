package enum

type PaymentMethod uint

const (
	PaymentMethodCash        PaymentMethod = iota + 1
	PaymentMethodInstallment
	PaymentMethodOnline
	PaymentMethodWallet
)

func (p PaymentMethod) String() string {
	switch p {
	case PaymentMethodCash:
		return "تفدی"
	case PaymentMethodInstallment:
		return "اقساط"
	case PaymentMethodOnline:
		return "آنلاین"
	case PaymentMethodWallet:
		return "کیف پول"
	}
	return "unknown"
}

func GetAllPaymentMethods() []PaymentMethod {
	return []PaymentMethod{
		PaymentMethodCash,
		PaymentMethodInstallment,
		PaymentMethodOnline,
		PaymentMethodWallet,
	}
}
