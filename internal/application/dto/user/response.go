package userdto

type OTPData struct {
	OTP      string `json:"otp"`
	Attempts int    `json:"attempts"`
}

type UserInfoResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Type         string `json:"type"`
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
}

type UserResponse struct {
	ID uint `json:"id"`
}

type AdminInfoResponse struct {
	AccessToken  string               `json:"accessToken"`
	RefreshToken string               `json:"refreshToken"`
}

type RevenueByTier struct {
	Tier     string `json:"tier"`
	Revenue  uint   `json:"revenue"`
}

type OrderByDay struct {
	Date  string `json:"date"`
	Count uint   `json:"count"`
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