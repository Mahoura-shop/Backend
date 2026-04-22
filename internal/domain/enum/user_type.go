package enum

type UserType uint

const (
	UserTypeGuest UserType = iota + 1
	UserTypeCustomer
	UserTypeShopkeeperCheque
	UserTypeShopkeeperCash
	UserTypeFellow
	UserTypeAdmin
)

func (userType UserType) String() string {
	switch userType {
	case UserTypeGuest:
		return "guest"
	case UserTypeCustomer:
		return "regular"
	case UserTypeShopkeeperCheque:
		return "shopkeeperCheque"
	case UserTypeShopkeeperCash:
		return "shopkeeperCash"
	case UserTypeFellow:
		return "fellow"
	case UserTypeAdmin:
		return "admin"
	}
	return ""
}

func GetAllUserTypes() []UserType {
	return []UserType{
		UserTypeGuest,
		UserTypeCustomer,
		UserTypeShopkeeperCheque,
		UserTypeShopkeeperCash,
		UserTypeFellow,
		UserTypeAdmin,
	}
}
