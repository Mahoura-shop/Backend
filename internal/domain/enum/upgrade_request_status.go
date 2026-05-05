package enum

type UpgradeRequestStatus uint

const (
	UpgradeRequestStatusPending       UpgradeRequestStatus = iota + 1
	UpgradeRequestStatusApproved
	UpgradeRequestStatusRejected
	UpgradeRequestStatusInfoRequested
)

func (s UpgradeRequestStatus) String() string {
	switch s {
	case UpgradeRequestStatusPending:
		return "pending"
	case UpgradeRequestStatusApproved:
		return "approved"
	case UpgradeRequestStatusRejected:
		return "rejected"
	case UpgradeRequestStatusInfoRequested:
		return "infoRequested"
	}
	return "unknown"
}
