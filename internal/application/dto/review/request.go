package reviewdto

type SubmitReviewRequest struct {
	UserID    uint
	ProductID uint
	Rating    uint   `json:"rating" validate:"required,min=1,max=5"`
	Comment   string `json:"comment"`
}
