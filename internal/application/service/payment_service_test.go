package service

import (
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/stretchr/testify/suite"
)

type PaymentServiceTestSuite struct {
	suite.Suite
	service *PaymentService
}

func (s *PaymentServiceTestSuite) SetupTest() {
	s.service = NewPaymentService(&bootstrap.Zarinpal{
		Enabled:     false,
		Sandbox:     true,
		MerchantID:  "test-merchant",
		CallbackURL: "https://example.com/callback",
	})
}

func (s *PaymentServiceTestSuite) TestIsEnabled_FalseWhenDisabled() {
	s.False(s.service.IsEnabled())
}

func (s *PaymentServiceTestSuite) TestIsEnabled_TrueWhenEnabled() {
	svc := NewPaymentService(&bootstrap.Zarinpal{Enabled: true})
	s.True(svc.IsEnabled())
}

func (s *PaymentServiceTestSuite) TestInitiateGatewayPayment_DisabledReturnsMockAuthority() {
	authority, gatewayURL, err := s.service.InitiateGatewayPayment(42, 1_000_000)

	s.NoError(err)
	s.Contains(authority, "mock_auth_42_1000000")
	s.Contains(gatewayURL, authority)
}

func (s *PaymentServiceTestSuite) TestInitiateGatewayPayment_DisabledURLContainsSandboxHost() {
	authority, gatewayURL, err := s.service.InitiateGatewayPayment(1, 500_000)

	s.NoError(err)
	s.NotEmpty(authority)
	s.Contains(gatewayURL, "sandbox.zarinpal.com")
}

func (s *PaymentServiceTestSuite) TestVerifyGatewayPayment_DisabledReturnsMockRefCode() {
	refCode, err := s.service.VerifyGatewayPayment("mock_auth_1_500000", 500_000)

	s.NoError(err)
	s.Equal("12345", refCode)
}

func (s *PaymentServiceTestSuite) TestInitiateGatewayPayment_ProductionURLWhenNotSandbox() {
	svc := NewPaymentService(&bootstrap.Zarinpal{
		Enabled: false,
		Sandbox: false,
	})

	_, gatewayURL, err := svc.InitiateGatewayPayment(1, 1_000_000)

	s.NoError(err)
	s.Contains(gatewayURL, "zarinpal.com")
	s.NotContains(gatewayURL, "sandbox")
}

func TestPaymentServiceSuite(t *testing.T) {
	suite.Run(t, new(PaymentServiceTestSuite))
}
