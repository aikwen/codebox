package models

import (
	"errors"
)
var ErrNoRecord = errors.New("models: no matching record found")

var ErrInvalidCredentials = errors.New("models: invalid credentials")

var ErrDuplicateEmail = errors.New("models: duplicate email")