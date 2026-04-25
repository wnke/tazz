package http

import (
	"errors"
)

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrDeviceNotFound = errors.New("device not found")
	ErrDeviceConflict = errors.New("device conflict")
)

type Power string

const (
	PowerOn      Power = "on"
	PowerOff     Power = "off"
	PowerStandby Power = "standby"
)

type Device struct {
	ID       string
	URL      string
	Username string
	Power    Power
}

type CreateDeviceInput struct {
	URL      string
	Username string
	Password string
	Power    Power
}

type UpdateDeviceInput struct {
	URL      *string
	Username *string
	Password *string
	Power    *Power
}

type ListDevicesFilter struct {
	Limit  int
	Offset int
}
