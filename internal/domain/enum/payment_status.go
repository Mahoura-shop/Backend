package enum

type PaymentStatus uint

const (
	PaymentStatusWaiting PaymentStatus = iota + 1
	PaymentStatusPaying
	PaymentStatusPaid
	PaymentStatusError
	PaymentStatusAll
)

func (cs PaymentStatus) String() string {
	switch cs {
	case PaymentStatusWaiting:
		return "در انتظار پرداخت"
	case PaymentStatusPaying:
		return "در حال پرداخت"
	case PaymentStatusPaid:
		return "پرداخت شده"
	case PaymentStatusError:
		return "خطا در پرداخت"
	case PaymentStatusAll:
		return "همه"
	}
	return "unknown"
}
func GetAllPaymentStatuses() []PaymentStatus {
	return []PaymentStatus{
		PaymentStatusWaiting,
		PaymentStatusPaying,
		PaymentStatusPaid,
		PaymentStatusError,
		PaymentStatusAll,
	}
}
