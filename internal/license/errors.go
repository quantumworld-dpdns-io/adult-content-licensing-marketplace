package license

import "errors"

var (
	ErrNotFound   = errors.New("license not found")
	ErrConflict   = errors.New("license version conflict")
	ErrBadVersion = errors.New("missing or invalid license version")
)
