package trip

import (
	"context"
	"database/sql"

	"github.com/go-logr/logr"
	"github.com/huydnt1801/chuyende/internal/ent"
	entTrip "github.com/huydnt1801/chuyende/internal/ent/trip"
	entutil "github.com/huydnt1801/chuyende/internal/utils/ent"
	"github.com/huydnt1801/chuyende/pkg/log"
)

type Service interface {
	FindTrip(ctx context.Context, params *TripParams) ([]*Trip, error)
	CreateTrip(ctx context.Context, trip *Trip) (*Trip, error)
	UpdateTrip(ctx context.Context, id int, updateParams *TripUpdate) (*Trip, error)
}

type ServiceImpl struct {
	logger    logr.Logger
	entClient *ent.Client
}

func NewService(db *sql.DB) *ServiceImpl {
	client := entutil.NewClientFromDB(db)
	s := &ServiceImpl{
		logger:    log.ZapLogger(),
		entClient: client,
	}
	return s
}

func (s *ServiceImpl) FindTrip(ctx context.Context, params *TripParams) ([]*Trip, error) {
	repo := NewRepo(s.entClient)
	return repo.FindTrip(ctx, params)
}

func (s *ServiceImpl) CreateTrip(ctx context.Context, trip *Trip) (*Trip, error) {
	repo := NewRepo(s.entClient)
	if trip.Type == entTrip.TypeMotor {
		trip.Price = 20 * 1000 * trip.Distance
	} else {
		trip.Price = 30 * 1000 * trip.Distance
	}
	return repo.CreateTrip(ctx, trip)
}

func (s *ServiceImpl) UpdateTrip(ctx context.Context, id int, updateParams *TripUpdate) (*Trip, error) {
	repo := NewRepo(s.entClient)
	return repo.UpdateTrip(ctx, id, updateParams)

}
