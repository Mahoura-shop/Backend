package communication

type SMSService interface {
	SendOTP(receptor string, token string) error
	SendMessage(receptor string, message string) error
}
