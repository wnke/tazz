package persistence

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/wnke/tazz/internal/adapter/driven/persistence/sqlc/gen"
	domain "github.com/wnke/tazz/internal/domain/adapter/driven/persistence"
)

type DeviceRepository struct {
	queries *db.Queries
}

func NewDeviceRepository(queries *db.Queries) *DeviceRepository {
	return &DeviceRepository{queries: queries}
}

func (r *DeviceRepository) Create(ctx context.Context, device domain.Device) (domain.Device, error) {
	created, err := r.queries.CreateDevice(ctx, db.CreateDeviceParams{
		URL:      device.URL,
		Username: device.Username,
		Password: device.Password,
	})
	if err != nil {
		return domain.Device{}, err
	}

	return domain.Device{
		ID:       created.ID,
		URL:      created.URL,
		Username: created.Username,
		Password: created.Password,
	}, nil
}

func (r *DeviceRepository) GetByID(ctx context.Context, id int64) (domain.Device, error) {
	device, err := r.queries.GetDeviceByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Device{}, domain.ErrDeviceNotFound
		}
		return domain.Device{}, err
	}

	return domain.Device{
		ID:       device.ID,
		URL:      device.URL,
		Username: device.Username,
		Password: device.Password,
	}, nil
}

func (r *DeviceRepository) List(ctx context.Context) ([]domain.Device, error) {
	devices, err := r.queries.ListDevices(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]domain.Device, 0, len(devices))
	for _, d := range devices {
		out = append(out, domain.Device{
			ID:       d.ID,
			URL:      d.URL,
			Username: d.Username,
			Password: d.Password,
		})
	}

	return out, nil
}

func (r *DeviceRepository) Update(ctx context.Context, device domain.Device) (domain.Device, error) {
	updated, err := r.queries.UpdateDevice(ctx, db.UpdateDeviceParams{
		URL:      device.URL,
		Username: device.Username,
		Password: device.Password,
		ID:       device.ID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Device{}, domain.ErrDeviceNotFound
		}
		return domain.Device{}, err
	}

	return domain.Device{
		ID:       updated.ID,
		URL:      updated.URL,
		Username: updated.Username,
		Password: updated.Password,
	}, nil
}

func (r *DeviceRepository) Delete(ctx context.Context, id int64) error {
	rowsAffected, err := r.queries.DeleteDevice(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrDeviceNotFound
	}

	return nil
}
