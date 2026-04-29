package usecase

type PaymentService interface {
	InitiateGatewayPayment(orderID uint, amount uint) (authority string, gatewayURL string, err error)
	VerifyGatewayPayment(authority string, amount uint) (refCode string, err error)
}
