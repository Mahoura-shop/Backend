package enum

type InstalmentStatus uint

const (
	InstalmentStatusPending InstalmentStatus = iota + 1
	InstalmentStatusPaid
	InstalmentStatusOverdue
)

func (s InstalmentStatus) String() string {
	switch s {
	case InstalmentStatusPending:
		return "در انتظار"
	case InstalmentStatusPaid:
		return "پرداخت شده"
	case InstalmentStatusOverdue:
		return "سررسید گذشته"
	}
	return "unknown"
}
