package upgraderequestdto

import "github.com/Mahoura-shop/Backend/internal/domain/enum"

type SubmitUpgradeRequestRequest struct {
	UserID        uint
	RequestedType enum.UserType `json:"requestedType" validate:"required"`
	BusinessName  string        `json:"businessName" validate:"required"`
	TaxID         string        `json:"taxID"`
}

type ReviewUpgradeRequestRequest struct {
	Action    string `json:"action" validate:"required,oneof=approve reject info"`
	AdminNote string `json:"adminNote"`
	AdminID   uint
}

type ChangeUserTypeRequest struct {
	NewType uint `json:"newType" validate:"required"`
	AdminID uint
}
