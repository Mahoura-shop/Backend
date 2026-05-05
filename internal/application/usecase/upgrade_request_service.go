package usecase

import upgraderequestdto "github.com/Mahoura-shop/Backend/internal/application/dto/upgrade_request"

type UpgradeRequestService interface {
	SubmitUpgradeRequest(upgraderequestdto.SubmitUpgradeRequestRequest) error
	GetMyUpgradeRequests(userID uint) ([]upgraderequestdto.UpgradeRequestCredential, error)
	GetAllUpgradeRequests(status string) ([]upgraderequestdto.UpgradeRequestCredential, error)
	GetUpgradeRequest(requestID uint) (*upgraderequestdto.UpgradeRequestCredential, error)
	ReviewUpgradeRequest(requestID uint, req upgraderequestdto.ReviewUpgradeRequestRequest) error
	ChangeUserType(userID uint, req upgraderequestdto.ChangeUserTypeRequest) error
	GetUserAuditLogs(userID uint) ([]upgraderequestdto.UserAuditLogCredential, error)
}
