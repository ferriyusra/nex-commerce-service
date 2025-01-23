package pagination

import "errors"

var (
	ErrorMaxPage      = errors.New("cooshen page more than total page")
	ErrorPage         = errors.New("page must greater than 0")
	ErrorPerPageEmpty = errors.New("page cannot be empty")
	ErrorPageInvalid  = errors.New("page invalid, must be number")
)
