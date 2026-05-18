package enum

type UserType uint

const (
	UserTypeGuest UserType = iota + 1
	UserTypeCustomer
	UserTypeShopkeeperCheque // legacy — treated same as ShopkeeperCash
	UserTypeShopkeeperCash
	UserTypeFellow
	UserTypeAdmin
)

// UserTypeShopkeeper is the canonical name going forward.
const UserTypeShopkeeper = UserTypeShopkeeperCash

func (userType UserType) String() string {
	switch userType {
	case UserTypeGuest:
		return "guest"
	case UserTypeCustomer:
		return "regular"
	case UserTypeShopkeeperCheque, UserTypeShopkeeperCash:
		return "shopkeeper"
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
		UserTypeShopkeeperCash,
		UserTypeFellow,
		UserTypeAdmin,
	}
}
