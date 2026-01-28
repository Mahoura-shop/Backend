package postgres

type TestRepository struct{}

func NewTestRepository() *TestRepository {
	return &TestRepository{}
}
