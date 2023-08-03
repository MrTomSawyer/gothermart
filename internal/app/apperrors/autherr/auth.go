package autherr

import "errors"

var ErrWrongCredentials = errors.New("wrong login or password")
var ErrJWTSignFailure = errors.New("failed to sign jwt")
