package usecase

type TestService interface {
	Test(string) (string, error)
}