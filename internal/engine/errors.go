package engine

import "errors"

var ErrProtectedPath = errors.New("path belongs to warren and is not writable by the model")
