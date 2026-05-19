package enum

type NotificationType uint

const (
	NotificationTypeOrderStatus NotificationType = iota + 1
	NotificationTypeBackInStock
)
