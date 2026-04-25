package http

import (
	"context"
	"errors"

	gen "github.com/wnke/tazz/internal/driver/http/gen"
)

type Service struct {
	manager interface {
		CreateDevice(ctx context.Context, in CreateDeviceInput) (Device, error)
		GetDevice(ctx context.Context, id string) (Device, error)
		ListDevices(ctx context.Context, filter ListDevicesFilter) ([]Device, error)
		UpdateDevice(ctx context.Context, id string, in UpdateDeviceInput) (Device, error)
		DeleteDevice(ctx context.Context, id string) error
	}
}

func NewService(manager interface {
	CreateDevice(ctx context.Context, in CreateDeviceInput) (Device, error)
	GetDevice(ctx context.Context, id string) (Device, error)
	ListDevices(ctx context.Context, filter ListDevicesFilter) ([]Device, error)
	UpdateDevice(ctx context.Context, id string, in UpdateDeviceInput) (Device, error)
	DeleteDevice(ctx context.Context, id string) error
}) *Service {
	return &Service{manager: manager}
}

func (s *Service) CreateDevice(ctx context.Context, opts *gen.CreateDeviceServiceRequestOptions) (*gen.CreateDeviceResponseData, error) {
	if opts == nil || opts.Body == nil || opts.Body.URL == "" || opts.Body.Username == "" || opts.Body.Password == nil || *opts.Body.Password == "" || opts.Body.Power == "" {
		return nil, gen.BadRequest{Code: "bad_request", Message: "missing required fields"}
	}

	in := CreateDeviceInput{
		URL:      opts.Body.URL,
		Username: opts.Body.Username,
		Password: *opts.Body.Password,
		Power:    Power(opts.Body.Power),
	}

	device, err := s.manager.CreateDevice(ctx, in)
	if err != nil {
		return nil, s.mapError(err)
	}

	resp := &gen.CreateDeviceResponse{
		ID:       device.ID,
		URL:      device.URL,
		Username: device.Username,
		Power:    gen.Power(device.Power),
	}
	return gen.NewCreateDeviceResponseData(resp), nil
}

func (s *Service) ListDevices(ctx context.Context, opts *gen.ListDevicesServiceRequestOptions) (*gen.ListDevicesResponseData, error) {
	filter := ListDevicesFilter{Limit: 20, Offset: 0}
	if opts != nil && opts.Query != nil {
		if opts.Query.Limit != nil {
			filter.Limit = *opts.Query.Limit
		}
		if opts.Query.Offset != nil {
			filter.Offset = *opts.Query.Offset
		}
	}
	if filter.Limit < 1 || filter.Limit > 100 || filter.Offset < 0 {
		return nil, gen.BadRequest{Code: "bad_request", Message: "invalid pagination"}
	}

	devices, err := s.manager.ListDevices(ctx, filter)
	if err != nil {
		return nil, s.mapError(err)
	}

	items := make([]gen.Device, 0, len(devices))
	for _, d := range devices {
		items = append(items, gen.Device{
			ID:       d.ID,
			URL:      d.URL,
			Username: d.Username,
			Power:    gen.Power(d.Power),
		})
	}

	resp := &gen.ListDevicesResponse{Items: items}
	return gen.NewListDevicesResponseData(resp), nil
}

func (s *Service) GetDevice(ctx context.Context, opts *gen.GetDeviceServiceRequestOptions) (*gen.GetDeviceResponseData, error) {
	if opts == nil || opts.PathParams == nil || opts.PathParams.ID == "" {
		return nil, gen.BadRequest{Code: "bad_request", Message: "id is required"}
	}

	device, err := s.manager.GetDevice(ctx, opts.PathParams.ID)
	if err != nil {
		return nil, s.mapError(err)
	}

	resp := &gen.GetDeviceResponse{
		ID:       device.ID,
		URL:      device.URL,
		Username: device.Username,
		Power:    gen.Power(device.Power),
	}
	return gen.NewGetDeviceResponseData(resp), nil
}

func (s *Service) UpdateDevice(ctx context.Context, opts *gen.UpdateDeviceServiceRequestOptions) (*gen.UpdateDeviceResponseData, error) {
	if opts == nil || opts.PathParams == nil || opts.PathParams.ID == "" {
		return nil, gen.BadRequest{Code: "bad_request", Message: "id is required"}
	}
	if opts.Body == nil {
		return nil, gen.BadRequest{Code: "bad_request", Message: "empty update payload"}
	}
	if opts.Body.URL == nil && opts.Body.Username == nil && opts.Body.Password == nil && opts.Body.Power == nil {
		return nil, gen.BadRequest{Code: "bad_request", Message: "empty update payload"}
	}

	in := UpdateDeviceInput{
		URL:      opts.Body.URL,
		Username: opts.Body.Username,
		Password: opts.Body.Password,
	}
	if opts.Body.Power != nil {
		power := Power(*opts.Body.Power)
		in.Power = &power
	}

	device, err := s.manager.UpdateDevice(ctx, opts.PathParams.ID, in)
	if err != nil {
		return nil, s.mapError(err)
	}

	resp := &gen.UpdateDeviceResponse{
		ID:       device.ID,
		URL:      device.URL,
		Username: device.Username,
		Power:    gen.Power(device.Power),
	}
	return gen.NewUpdateDeviceResponseData(resp), nil
}

func (s *Service) DeleteDevice(ctx context.Context, opts *gen.DeleteDeviceServiceRequestOptions) (*gen.DeleteDeviceResponseData, error) {
	if opts == nil || opts.PathParams == nil || opts.PathParams.ID == "" {
		return nil, gen.BadRequest{Code: "bad_request", Message: "id is required"}
	}

	if err := s.manager.DeleteDevice(ctx, opts.PathParams.ID); err != nil {
		return nil, s.mapError(err)
	}

	return gen.NewDeleteDeviceResponseData(nil), nil
}

func (s *Service) mapError(err error) error {
	switch {
	case errors.Is(err, ErrInvalidInput):
		return gen.BadRequest{Code: "bad_request", Message: err.Error()}
	case errors.Is(err, ErrDeviceNotFound):
		return gen.NotFound{Code: "not_found", Message: err.Error()}
	case errors.Is(err, ErrDeviceConflict):
		return gen.Conflict{Code: "conflict", Message: err.Error()}
	default:
		return gen.InternalError{Code: "internal_error", Message: "internal server error"}
	}
}
