package upgraderequestdto

import (
	"time"

	"github.com/Mahoura-shop/Backend/internal/domain/enum"
)

type UpgradeRequestCredential struct {
	ID            uint                      `json:"id"`
	UserID        uint                      `json:"userID"`
	UserPhone     string                    `json:"userPhone"`
	RequestedType string                    `json:"requestedType"`
	BusinessName  string                    `json:"businessName"`
	TaxID         string                    `json:"taxID"`
	Status        enum.UpgradeRequestStatus `json:"status"`
	StatusLabel   string                    `json:"statusLabel"`
	AdminNote     string                    `json:"adminNote"`
	ReviewedByID  *uint                     `json:"reviewedByID"`
	CreatedAt     time.Time                 `json:"createdAt"`
}

type UserAuditLogCredential struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"userID"`
	ChangedByID uint      `json:"changedByID"`
	OldType     string    `json:"oldType"`
	NewType     string    `json:"newType"`
	Reason      string    `json:"reason"`
	CreatedAt   time.Time `json:"createdAt"`
}
