// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/huydnt1801/chuyende/internal/ent/otp"
)

// Otp is the model entity for the Otp schema.
type Otp struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// PhoneNumber holds the value of the "phone_number" field.
	PhoneNumber string `json:"phone_number,omitempty"`
	// Otp holds the value of the "otp" field.
	Otp string `json:"otp,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt    time.Time `json:"created_at,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Otp) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case otp.FieldID:
			values[i] = new(sql.NullInt64)
		case otp.FieldPhoneNumber, otp.FieldOtp:
			values[i] = new(sql.NullString)
		case otp.FieldCreatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Otp fields.
func (o *Otp) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case otp.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			o.ID = int(value.Int64)
		case otp.FieldPhoneNumber:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field phone_number", values[i])
			} else if value.Valid {
				o.PhoneNumber = value.String
			}
		case otp.FieldOtp:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field otp", values[i])
			} else if value.Valid {
				o.Otp = value.String
			}
		case otp.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				o.CreatedAt = value.Time
			}
		default:
			o.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Otp.
// This includes values selected through modifiers, order, etc.
func (o *Otp) Value(name string) (ent.Value, error) {
	return o.selectValues.Get(name)
}

// Update returns a builder for updating this Otp.
// Note that you need to call Otp.Unwrap() before calling this method if this Otp
// was returned from a transaction, and the transaction was committed or rolled back.
func (o *Otp) Update() *OtpUpdateOne {
	return NewOtpClient(o.config).UpdateOne(o)
}

// Unwrap unwraps the Otp entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (o *Otp) Unwrap() *Otp {
	_tx, ok := o.config.driver.(*txDriver)
	if !ok {
		panic("ent: Otp is not a transactional entity")
	}
	o.config.driver = _tx.drv
	return o
}

// String implements the fmt.Stringer.
func (o *Otp) String() string {
	var builder strings.Builder
	builder.WriteString("Otp(")
	builder.WriteString(fmt.Sprintf("id=%v, ", o.ID))
	builder.WriteString("phone_number=")
	builder.WriteString(o.PhoneNumber)
	builder.WriteString(", ")
	builder.WriteString("otp=")
	builder.WriteString(o.Otp)
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(o.CreatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// Otps is a parsable slice of Otp.
type Otps []*Otp