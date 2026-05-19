package enum

type OrderStatus uint

const (
	OrderStatusPending   OrderStatus = 1
	OrderStatusPaid      OrderStatus = 2
	OrderStatusShipped   OrderStatus = 3
	OrderStatusCancelled OrderStatus = 5
)

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusPending:
		return "در انتظار پرداخت"
	case OrderStatusPaid:
		return "پرداخت شده"
	case OrderStatusShipped:
		return "ارسال شده"
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
		OrderStatusCancelled,
	}
}
