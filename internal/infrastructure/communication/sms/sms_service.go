package sms

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/kavenegar/kavenegar-go"
)

type SMSService struct {
	providerConfig *bootstrap.SMSGateway
	smsTemplates   *bootstrap.SMSTemplates
}

func NewSMSService(
	providerConfig *bootstrap.SMSGateway,
	smsTemplates *bootstrap.SMSTemplates,
) *SMSService {
	return &SMSService{
		providerConfig: providerConfig,
		smsTemplates:   smsTemplates,
	}
}

func (smsService *SMSService) SendOTP(receptor, token string) error {
	if !smsService.providerConfig.Enabled {
		return nil
	}
	api := kavenegar.New(smsService.providerConfig.APIKey)
	template := smsService.smsTemplates.OTP
	params := &kavenegar.VerifyLookupParam{}
	if _, err := api.Verify.Lookup(receptor, template, token, params); err != nil {
		return err
	}
	return nil
}

func (smsService *SMSService) SendMessage(receptor, message string) error {
	if !smsService.providerConfig.Enabled {
		return nil
	}
	api := kavenegar.New(smsService.providerConfig.APIKey)
	_, err := api.Message.Send(smsService.providerConfig.Sender, []string{receptor}, message, nil)
	return err
}
