package enum

type CartStatus uint

const (
	CartStatusWaiting CartStatus = iota + 1
	CartStatusConfirmed
	CartStatusAll
)

func (cs CartStatus) String() string {
	switch cs {
	case CartStatusWaiting:
		return "در انتظار تایید"
	case CartStatusConfirmed:
		return "تایید شده"
	case CartStatusAll:
		return "همه"
	}
	return "unknown"
}
func GetAllCartStatuses() []CartStatus {
	return []CartStatus{
		CartStatusWaiting,
		CartStatusConfirmed,
		CartStatusAll,
	}
}
