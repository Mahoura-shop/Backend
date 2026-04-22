package enum

type TransactionType uint

const (
	TransactionTypeDeposit TransactionType = iota + 1
	TransactionTypeWithDraw
	TransactionTypeAll
)

func (userType TransactionType) String() string {
	switch userType {
	case TransactionTypeDeposit:
		return "deposit"
	case TransactionTypeWithDraw:
		return "withdraw"
	case TransactionTypeAll:
		return "all"
	}
	return ""
}

func GetAllTransactionTypes() []TransactionType {
	return []TransactionType{
		TransactionTypeDeposit,
		TransactionTypeWithDraw,
		TransactionTypeAll,
	}
}
