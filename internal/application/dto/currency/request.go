package currencydto

type CreateCurrencyRequest struct {
	Name        string
	Code        string
	ConvertRate uint
}

type UpdateCurrencyRequest struct {
	ID          uint
	Name        *string
	Code        *string
	ConvertRate *uint
}