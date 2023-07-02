package trip

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/huydnt1801/chuyende/internal/ent"
	"github.com/huydnt1801/chuyende/internal/ent/trip"
	"github.com/huydnt1801/chuyende/pkg/log"
)

type Repo interface {
	FindTrip(ctx context.Context, params *TripParams) ([]*Trip, error)
	CreateTrip(ctx context.Context, Trip *Trip) (*Trip, error)
	UpdateTrip(ctx context.Context, id int, updated *TripUpdate) (*Trip, error)
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

func (r *RepoImpl) FindTrip(ctx context.Context, params *TripParams) ([]*Trip, error) {
	q := r.client.Trip.Query().WithUser().WithDriver()
	if v := params.TripID; v != nil {
		q = q.Where(trip.ID(*v))
	}
	if v := params.UserID; v != nil {
		q = q.Where(trip.UserID(*v))
	}
	if v := params.DriveID; v != nil {
		q = q.Where(trip.DriverID(*v))
	}
	if v := params.Status; v != nil {
		q = q.Where(trip.StatusEQ(*v))
	}
	if v := params.Rate; v != nil {
		q = q.Where(trip.Rate(*v))
	}

	u, err := q.All(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, TripNotFoundError{}
		}
		return nil, err
	}
	return DecodeManyTrip(u)
}

func (r *RepoImpl) CreateTrip(ctx context.Context, trip *Trip) (*Trip, error) {
	q := r.client.Trip.Create().
		SetUserID(trip.UserID).
		SetStartX(trip.StartX).
		SetStartY(trip.StartY).
		SetStartLocation(trip.StartLocation).
		SetEndX(trip.EndX).
		SetEndY(trip.EndY).
		SetEndLocation(trip.EndLocation).
		SetType(trip.Type).
		SetPrice(trip.Price).
		SetDistance(trip.Distance)
	u, err := q.Save(ctx)
	if ent.IsValidationError(err) {
		return nil, InvalidTripError{}
	}
	if err != nil {
		return nil, fmt.Errorf("failed creating Trip: %w", err)
	}
	return DecodeTrip(u)
}

func (r *RepoImpl) UpdateTrip(ctx context.Context, id int, updated *TripUpdate) (*Trip, error) {
	q := r.client.Trip.UpdateOneID(id)

	if v := updated.DriveID; v != nil {
		q.SetDriverID(*v)
	}
	if v := updated.Status; v != nil {
		q.SetStatus(*v)
	}
	if v := updated.Rate; v != nil {
		q.SetRate(*v)
	}
	trip, err := q.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed updating Trip: %w", err)
	}
	return DecodeTrip(trip)
}
