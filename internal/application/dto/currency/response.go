package currencydto

type CurrencyCredential struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	ConvertRate uint   `json:"convertRate"`
}