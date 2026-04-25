package http

import (
	"context"
	"net/http"

	gen "github.com/wnke/tazz/internal/driver/http/gen"
)

type DeviceManager interface {
	CreateDevice(ctx context.Context, in CreateDeviceInput) (Device, error)
	GetDevice(ctx context.Context, id string) (Device, error)
	ListDevices(ctx context.Context, filter ListDevicesFilter) ([]Device, error)
	UpdateDevice(ctx context.Context, id string, in UpdateDeviceInput) (Device, error)
	DeleteDevice(ctx context.Context, id string) error
}

func New(manager DeviceManager) http.Handler {
	svc := NewService(manager)
	return gen.NewRouter(svc)
}
