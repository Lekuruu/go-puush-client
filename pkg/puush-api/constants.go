package puush

type AccountType int8

const (
	AccountTypeRegular AccountType = iota
	AccountTypePro
	AccountTypeUnlimited
)
