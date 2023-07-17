package withdrawerr

import "errors"

var ErrLowBalance = errors.New("Not enough points to withdraw")
