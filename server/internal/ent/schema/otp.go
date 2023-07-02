package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Otp holds the schema definition for the Otp entity.
type Otp struct {
	ent.Schema
}

// Fields of the Otp.
func (Otp) Fields() []ent.Field {
	return []ent.Field{
		field.String("phone_number").Unique(),
		field.String("otp"),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Otp.
func (Otp) Edges() []ent.Edge {
	return nil
}
