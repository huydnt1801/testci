// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/huydnt1801/chuyende/internal/ent/vehicledriver"
)

// VehicleDriver is the model entity for the VehicleDriver schema.
type VehicleDriver struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// PhoneNumber holds the value of the "phone_number" field.
	PhoneNumber string `json:"phone_number,omitempty"`
	// FullName holds the value of the "full_name" field.
	FullName string `json:"full_name,omitempty"`
	// Password holds the value of the "password" field.
	Password string `json:"password,omitempty"`
	// License holds the value of the "license" field.
	License vehicledriver.License `json:"license,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the VehicleDriverQuery when eager-loading is set.
	Edges        VehicleDriverEdges `json:"edges"`
	selectValues sql.SelectValues
}

// VehicleDriverEdges holds the relations/edges for other nodes in the graph.
type VehicleDriverEdges struct {
	// Trips holds the value of the trips edge.
	Trips []*Trip `json:"trips,omitempty"`
	// Sessions holds the value of the sessions edge.
	Sessions []*Session `json:"sessions,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// TripsOrErr returns the Trips value or an error if the edge
// was not loaded in eager-loading.
func (e VehicleDriverEdges) TripsOrErr() ([]*Trip, error) {
	if e.loadedTypes[0] {
		return e.Trips, nil
	}
	return nil, &NotLoadedError{edge: "trips"}
}

// SessionsOrErr returns the Sessions value or an error if the edge
// was not loaded in eager-loading.
func (e VehicleDriverEdges) SessionsOrErr() ([]*Session, error) {
	if e.loadedTypes[1] {
		return e.Sessions, nil
	}
	return nil, &NotLoadedError{edge: "sessions"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*VehicleDriver) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case vehicledriver.FieldID:
			values[i] = new(sql.NullInt64)
		case vehicledriver.FieldPhoneNumber, vehicledriver.FieldFullName, vehicledriver.FieldPassword, vehicledriver.FieldLicense:
			values[i] = new(sql.NullString)
		case vehicledriver.FieldCreatedAt, vehicledriver.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the VehicleDriver fields.
func (vd *VehicleDriver) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case vehicledriver.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			vd.ID = int(value.Int64)
		case vehicledriver.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				vd.CreatedAt = value.Time
			}
		case vehicledriver.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				vd.UpdatedAt = value.Time
			}
		case vehicledriver.FieldPhoneNumber:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field phone_number", values[i])
			} else if value.Valid {
				vd.PhoneNumber = value.String
			}
		case vehicledriver.FieldFullName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field full_name", values[i])
			} else if value.Valid {
				vd.FullName = value.String
			}
		case vehicledriver.FieldPassword:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field password", values[i])
			} else if value.Valid {
				vd.Password = value.String
			}
		case vehicledriver.FieldLicense:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field license", values[i])
			} else if value.Valid {
				vd.License = vehicledriver.License(value.String)
			}
		default:
			vd.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the VehicleDriver.
// This includes values selected through modifiers, order, etc.
func (vd *VehicleDriver) Value(name string) (ent.Value, error) {
	return vd.selectValues.Get(name)
}

// QueryTrips queries the "trips" edge of the VehicleDriver entity.
func (vd *VehicleDriver) QueryTrips() *TripQuery {
	return NewVehicleDriverClient(vd.config).QueryTrips(vd)
}

// QuerySessions queries the "sessions" edge of the VehicleDriver entity.
func (vd *VehicleDriver) QuerySessions() *SessionQuery {
	return NewVehicleDriverClient(vd.config).QuerySessions(vd)
}

// Update returns a builder for updating this VehicleDriver.
// Note that you need to call VehicleDriver.Unwrap() before calling this method if this VehicleDriver
// was returned from a transaction, and the transaction was committed or rolled back.
func (vd *VehicleDriver) Update() *VehicleDriverUpdateOne {
	return NewVehicleDriverClient(vd.config).UpdateOne(vd)
}

// Unwrap unwraps the VehicleDriver entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (vd *VehicleDriver) Unwrap() *VehicleDriver {
	_tx, ok := vd.config.driver.(*txDriver)
	if !ok {
		panic("ent: VehicleDriver is not a transactional entity")
	}
	vd.config.driver = _tx.drv
	return vd
}

// String implements the fmt.Stringer.
func (vd *VehicleDriver) String() string {
	var builder strings.Builder
	builder.WriteString("VehicleDriver(")
	builder.WriteString(fmt.Sprintf("id=%v, ", vd.ID))
	builder.WriteString("created_at=")
	builder.WriteString(vd.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(vd.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("phone_number=")
	builder.WriteString(vd.PhoneNumber)
	builder.WriteString(", ")
	builder.WriteString("full_name=")
	builder.WriteString(vd.FullName)
	builder.WriteString(", ")
	builder.WriteString("password=")
	builder.WriteString(vd.Password)
	builder.WriteString(", ")
	builder.WriteString("license=")
	builder.WriteString(fmt.Sprintf("%v", vd.License))
	builder.WriteByte(')')
	return builder.String()
}

// VehicleDrivers is a parsable slice of VehicleDriver.
type VehicleDrivers []*VehicleDriver
