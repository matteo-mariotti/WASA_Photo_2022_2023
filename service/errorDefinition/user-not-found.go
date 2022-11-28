package errorDefinition

import "errors"

var ErrUserNotFound = errors.New("User not found in the database")
