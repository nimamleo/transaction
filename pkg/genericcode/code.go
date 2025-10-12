package genericcode

type Code int

const (
	InternalServerError Code = iota + 1
	NotFound
	Unauthorized
	OK
	Forbidden
	BadRequest
	Conflict
)
