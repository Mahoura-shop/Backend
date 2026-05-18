package usecase

import returndto "github.com/Mahoura-shop/Backend/internal/application/dto/return"

type ReturnService interface {
	RequestReturn(req returndto.RequestReturnRequest) error
	GetMyReturns(userID uint) ([]returndto.ReturnCredential, error)
	GetAllReturns(status string) ([]returndto.ReturnCredential, error)
	ReviewReturn(returnID uint, req returndto.ReviewReturnRequest) error
	ProcessRefund(returnID uint, adminID uint) error
}
