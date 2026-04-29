package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mahoura-shop/Backend/bootstrap"
)

type PaymentService struct {
	config *bootstrap.Zarinpal
}

func NewPaymentService(config *bootstrap.Zarinpal) *PaymentService {
	return &PaymentService{config: config}
}

func (s *PaymentService) baseURL() string {
	if s.config.Sandbox {
		return "https://sandbox.zarinpal.com/pg/v4/payment"
	}
	return "https://api.zarinpal.com/pg/v4/payment"
}

func (s *PaymentService) gatewayHost() string {
	if s.config.Sandbox {
		return "https://sandbox.zarinpal.com/pg/StartPay"
	}
	return "https://www.zarinpal.com/pg/StartPay"
}

func (s *PaymentService) InitiateGatewayPayment(orderID uint, amount uint) (string, string, error) {
	payload := map[string]interface{}{
		"merchant_id":  s.config.MerchantID,
		"amount":       amount,
		"callback_url": fmt.Sprintf("%s?orderID=%d", s.config.CallbackURL, orderID),
		"description":  fmt.Sprintf("پرداخت سفارش شماره %d", orderID),
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(s.baseURL()+"/request.json", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Code      int    `json:"code"`
			Authority string `json:"authority"`
		} `json:"data"`
		Errors interface{} `json:"errors"`
	}

	respBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", "", err
	}

	if result.Data.Code != 100 {
		return "", "", fmt.Errorf("zarinpal error code: %d", result.Data.Code)
	}

	authority := result.Data.Authority
	gatewayURL := fmt.Sprintf("%s/%s", s.gatewayHost(), authority)
	return authority, gatewayURL, nil
}

func (s *PaymentService) VerifyGatewayPayment(authority string, amount uint) (string, error) {
	payload := map[string]interface{}{
		"merchant_id": s.config.MerchantID,
		"amount":      amount,
		"authority":   authority,
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(s.baseURL()+"/verify.json", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Code    int    `json:"code"`
			RefID   int64  `json:"ref_id"`
			Message string `json:"message"`
		} `json:"data"`
	}

	respBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	if result.Data.Code != 100 && result.Data.Code != 101 {
		return "", fmt.Errorf("zarinpal verification failed: %s", result.Data.Message)
	}

	return fmt.Sprintf("%d", result.Data.RefID), nil
}
