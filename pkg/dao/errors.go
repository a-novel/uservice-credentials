package dao

import "errors"

var ErrCredentialsNotFound = errors.New("credentials not found")

var ErrCredentialsAlreadyExist = errors.New("credentials already exist")
