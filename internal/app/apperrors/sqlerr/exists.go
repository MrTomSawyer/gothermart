package sqlerr

import "errors"

var ErrUploadedBySameUser = errors.New("already uploaded by the same user")
var ErrUploadedByAnotherUser = errors.New("already uploaded by another user")
