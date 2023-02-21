package localErrors

import "errors"

var (
	ErrNotFoundChildren = errors.New("not found children links")
)
