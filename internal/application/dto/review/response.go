package reviewdto

type ReviewCredential struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"userID"`
	UserPhone     string `json:"userPhone"`
	UserFirstName string `json:"userFirstName"`
	UserLastName  string `json:"userLastName"`
	ProductID     uint   `json:"productID"`
	ProductName   string `json:"productName"`
	ProductSlug   string `json:"productSlug"`
	ProductPic    string `json:"productPic"`
	Rating        uint   `json:"rating"`
	Comment       string `json:"comment"`
	IsVerified    bool   `json:"isVerified"`
	CreatedAt     string `json:"createdAt"`
}
