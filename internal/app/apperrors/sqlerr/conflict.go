package sqlerr

import "errors"

var ErrLoginConflict = errors.New("login already exists")
