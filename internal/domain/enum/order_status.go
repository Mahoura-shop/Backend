package enum

type OrderStatus uint

const (
	OrderStatusPending   OrderStatus = iota + 1
	OrderStatusPaid
	OrderStatusShipped
	OrderStatusDelivered
	OrderStatusCancelled
)

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusPending:
		return "در انتظار پرداخت"
	case OrderStatusPaid:
		return "پرداخت شده"
	case OrderStatusShipped:
		return "ارسال شده"
	case OrderStatusDelivered:
		return "تحویل داده شده"
	case OrderStatusCancelled:
		return "لغو شده"
	}
	return "unknown"
}

func GetAllOrderStatuses() []OrderStatus {
	return []OrderStatus{
		OrderStatusPending,
		OrderStatusPaid,
		OrderStatusShipped,
		OrderStatusDelivered,
		OrderStatusCancelled,
	}
}
