package errors

type ErrorCode int

func (e ErrorCode) Code() int {
	return int(e)
}

const (
	// common 1 -> 1000
	System ErrorCode = iota + 1
	Validation
	BadRequest
	NotAuthorized
	Forbidden
	NotFound
	AlreadyExist
	// module specific 1001 -> 2000
)
