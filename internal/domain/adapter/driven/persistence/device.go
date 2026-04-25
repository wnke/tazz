package domain

import "errors"

var (
	ErrDeviceNotFound = errors.New("device not found")
)

type Device struct {
	ID       int64
	URL      string
	Username string
	Password string
}
