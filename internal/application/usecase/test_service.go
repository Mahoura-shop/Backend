package usecase

type TestService interface {
	Test(test string) (string, error)
}