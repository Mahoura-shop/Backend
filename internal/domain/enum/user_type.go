package enum

type UserType uint

const (
	UserTypeCustomer         UserType = 2
	UserTypeShopkeeperCheque UserType = 3 // legacy — treated same as ShopkeeperCash
	UserTypeShopkeeperCash   UserType = 4
	UserTypeFellow           UserType = 5
	UserTypeAdmin            UserType = 6
)

// UserTypeShopkeeper is the canonical name going forward.
const UserTypeShopkeeper = UserTypeShopkeeperCash

func (userType UserType) String() string {
	switch userType {
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
		UserTypeCustomer,
		UserTypeShopkeeperCash,
		UserTypeFellow,
		UserTypeAdmin,
	}
}
