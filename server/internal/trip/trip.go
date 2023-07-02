package trip

import (
	"time"

	"github.com/huydnt1801/chuyende/internal/driver"
	"github.com/huydnt1801/chuyende/internal/ent/trip"
	"github.com/huydnt1801/chuyende/internal/user"
)

type Trip struct {
	ID            int            `json:"id"`
	User          *user.User     `json:"user"`
	Drive         *driver.Driver `json:"drive,omitempty"`
	UserID        int            `json:"userId"`
	DriveID       int            `json:"driverId,omitempty" mapstructure:"driverId,omitempty"`
	StartX        float64        `json:"startX"`
	StartY        float64        `json:"startY"`
	StartLocation string         `json:"startLocation"`
	EndX          float64        `json:"endX"`
	EndY          float64        `json:"endY"`
	EndLocation   string         `json:"endLocation"`
	Distance      float64        `json:"distance"`
	Type          trip.Type      `json:"type"`
	Price         float64        `json:"price"`
	Status        trip.Status    `json:"status,omitempty"`
	Rate          int            `json:"rate,omitempty" mapstructure:"rate,omitempty"`
	CreatedAt     time.Time      `json:"createdAt"`
}

type TripParams struct {
	TripID  *int
	UserID  *int
	DriveID *int
	Status  *trip.Status
	Rate    *int
}

type TripUpdate struct {
	DriveID *int
	Status  *trip.Status
	Rate    *int
}
