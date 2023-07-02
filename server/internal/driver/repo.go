package driver

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/huydnt1801/chuyende/internal/ent"
	"github.com/huydnt1801/chuyende/internal/ent/vehicledriver"
	"github.com/huydnt1801/chuyende/pkg/log"
)

type Repo interface {
	FindDriver(ctx context.Context, phoneNumber string) (*Driver, error)
}

type RepoImpl struct {
	logger logr.Logger
	client *ent.Client
}

func NewRepo(client *ent.Client) *RepoImpl {
	return &RepoImpl{
		logger: log.ZapLogger(),
		client: client,
	}
}

func (r *RepoImpl) FindDriver(ctx context.Context, phoneNumber string) (*Driver, error) {
	q := r.client.VehicleDriver.Query()
	q = q.Where(vehicledriver.PhoneNumber(phoneNumber))

	u, err := q.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, DriverNotFoundError{}
		}
		return nil, err
	}
	return DecodeDriver(u)
}
