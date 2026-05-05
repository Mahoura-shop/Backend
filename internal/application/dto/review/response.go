package reviewdto

type ReviewCredential struct {
	ID         uint   `json:"id"`
	UserID     uint   `json:"userID"`
	UserPhone  string `json:"userPhone"`
	ProductID  uint   `json:"productID"`
	Rating     uint   `json:"rating"`
	Comment    string `json:"comment"`
	IsVerified bool   `json:"isVerified"`
	CreatedAt  string `json:"createdAt"`
}
