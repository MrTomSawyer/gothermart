package withdrawerr

import "errors"

var ErrLowBalance = errors.New("not enough points to withdraw")
