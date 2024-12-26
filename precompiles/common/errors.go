package common

const (
	ErrGetStateDB          = "get EVM StateDB failed"
	ErrInvalidNumberOfArgs = "invalid number of arguments; expected %d; got: %d"
	ErrSenderNotOrigin     = "msg.sender is not from tx origin"
	ErrWriteOnReadOnly     = "read only call to write functions"
)
