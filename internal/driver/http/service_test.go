package http

import (
	"context"
	"errors"
	"testing"

	gen "github.com/wnke/tazz/internal/driver/http/gen"
)

type stubManager struct {
	createFn func(ctx context.Context, in CreateDeviceInput) (Device, error)
	getFn    func(ctx context.Context, id string) (Device, error)
	listFn   func(ctx context.Context, filter ListDevicesFilter) ([]Device, error)
	updateFn func(ctx context.Context, id string, in UpdateDeviceInput) (Device, error)
	deleteFn func(ctx context.Context, id string) error
}

func (s stubManager) CreateDevice(ctx context.Context, in CreateDeviceInput) (Device, error) {
	return s.createFn(ctx, in)
}
func (s stubManager) GetDevice(ctx context.Context, id string) (Device, error) {
	return s.getFn(ctx, id)
}
func (s stubManager) ListDevices(ctx context.Context, filter ListDevicesFilter) ([]Device, error) {
	return s.listFn(ctx, filter)
}
func (s stubManager) UpdateDevice(ctx context.Context, id string, in UpdateDeviceInput) (Device, error) {
	return s.updateFn(ctx, id, in)
}
func (s stubManager) DeleteDevice(ctx context.Context, id string) error {
	return s.deleteFn(ctx, id)
}

func TestCreateDeviceSuccess(t *testing.T) {
	pwd := "secret"
	svc := NewService(stubManager{createFn: func(ctx context.Context, in CreateDeviceInput) (Device, error) {
		return Device{ID: "d-1", URL: in.URL, Username: in.Username, Power: in.Power}, nil
	}})

	resp, err := svc.CreateDevice(context.Background(), &gen.CreateDeviceServiceRequestOptions{Body: &gen.CreateDeviceBody{
		URL:      "https://example.com",
		Username: "u",
		Password: &pwd,
		Power:    gen.On,
	}})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp == nil || resp.Body == nil {
		t.Fatalf("expected response body")
	}
	if resp.Body.ID != "d-1" {
		t.Fatalf("expected id d-1 got %q", resp.Body.ID)
	}
}

func TestGetDeviceNotFound(t *testing.T) {
	svc := NewService(stubManager{getFn: func(ctx context.Context, id string) (Device, error) {
		return Device{}, ErrDeviceNotFound
	}})

	_, err := svc.GetDevice(context.Background(), &gen.GetDeviceServiceRequestOptions{PathParams: &gen.GetDevicePath{ID: "missing"}})
	if err == nil {
		t.Fatalf("expected error")
	}
	var notFound gen.NotFound
	if !errors.As(err, &notFound) {
		t.Fatalf("expected NotFound, got %T", err)
	}
}

func TestListDevicesSuccess(t *testing.T) {
	svc := NewService(stubManager{listFn: func(ctx context.Context, filter ListDevicesFilter) ([]Device, error) {
		return []Device{{ID: "d-1", URL: "https://a", Username: "u", Power: PowerOff}}, nil
	}})

	resp, err := svc.ListDevices(context.Background(), &gen.ListDevicesServiceRequestOptions{Query: &gen.ListDevicesQuery{}})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp == nil || resp.Body == nil || len(resp.Body.Items) != 1 {
		t.Fatalf("expected one device")
	}
}

func TestUpdateDeviceEmptyPayload(t *testing.T) {
	svc := NewService(stubManager{updateFn: func(ctx context.Context, id string, in UpdateDeviceInput) (Device, error) {
		return Device{}, nil
	}})

	_, err := svc.UpdateDevice(context.Background(), &gen.UpdateDeviceServiceRequestOptions{PathParams: &gen.UpdateDevicePath{ID: "d-1"}, Body: &gen.UpdateDeviceBody{}})
	if err == nil {
		t.Fatalf("expected error")
	}
	var badRequest gen.BadRequest
	if !errors.As(err, &badRequest) {
		t.Fatalf("expected BadRequest, got %T", err)
	}
}

func TestDeleteDeviceInternal(t *testing.T) {
	svc := NewService(stubManager{deleteFn: func(ctx context.Context, id string) error {
		return errors.New("boom")
	}})

	_, err := svc.DeleteDevice(context.Background(), &gen.DeleteDeviceServiceRequestOptions{PathParams: &gen.DeleteDevicePath{ID: "d-1"}})
	if err == nil {
		t.Fatalf("expected error")
	}
	var internal gen.InternalError
	if !errors.As(err, &internal) {
		t.Fatalf("expected InternalError, got %T", err)
	}
}
