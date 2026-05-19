package userdto

type OTPData struct {
	OTP      string `json:"otp"`
	Attempts int    `json:"attempts"`
}

type UserInfoResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	Type         string   `json:"type"`
	IsAdmin      bool     `json:"isAdmin"`
	Permissions  []string `json:"permissions"`
}

type UserCredential struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	ProfilePic string `json:"profilePic"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	IsAdmin    bool   `json:"isAdmin"`
	RoleID     *uint  `json:"roleID"`
	RoleName   string `json:"roleName"`
	CreatedAt  string `json:"createdAt"`
}

type UserResponse struct {
	ID uint `json:"id"`
}

type AdminInfoResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

type RevenueByTier struct {
	Tier     string `json:"tier"`
	Revenue  uint   `json:"revenue"`
}

type OrderByDay struct {
	Date  string `json:"date"`
	Count uint   `json:"count"`
}

type RevenueByDay struct {
	Date    string `json:"date"`
	Revenue uint   `json:"revenue"`
}

type LowStockProduct struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Quantity  uint   `json:"quantity"`
	MinOrder  uint   `json:"minOrder"`
}

type TopProduct struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Quantity uint   `json:"quantity"`
	Revenue  uint   `json:"revenue"`
}

type DashboardResponse struct {
	ProductsCount     uint              `json:"productsCount"`
	CategoriesCount   uint              `json:"categoriesCount"`
	BrandsCount       uint              `json:"brandsCount"`
	RevenuePerTier    []RevenueByTier   `json:"revenuePerTier"`
	OrdersPerDay      []OrderByDay      `json:"ordersPerDay"`
	LowStockProducts  []LowStockProduct `json:"lowStockProducts"`
	TopProducts       []TopProduct      `json:"topProducts"`
}

type UserWalletBalance struct {
	Balance uint `json:"balance"`
}

type TransactionDTO struct {
	ID        uint   `json:"id"`
	Amount    uint   `json:"amount"`
	Type      uint   `json:"type"`
	CreatedAt string `json:"createdAt"`
}

type AdminUserWalletResponse struct {
	Balance      uint             `json:"balance"`
	Transactions []TransactionDTO `json:"transactions"`
}

type SubAdminCredential struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	RoleID    *uint  `json:"roleID"`
	RoleName  string `json:"roleName"`
}