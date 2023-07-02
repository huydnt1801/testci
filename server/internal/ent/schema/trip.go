package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Trip holds the schema definition for the Trip entity.
type Trip struct {
	ent.Schema
}

func (Trip) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Trip.
func (Trip) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id"),
		field.Int("driver_id").
			Optional().
			StructTag(`mapstructure:",omitempty"`),
		field.Float("start_x"),
		field.Float("start_y"),
		field.String("start_location"),
		field.Float("end_x"),
		field.Float("end_y"),
		field.String("end_location"),
		field.Float("distance"),
		field.Float("price"),
		field.Enum("type").
			Values("motor", "car"),
		field.Enum("status").
			Values("waiting", "accept", "done", "cancel").
			Default("waiting"),
		field.Int("rate").
			Min(1).
			Max(5).
			Optional().
			StructTag(`mapstructure:",omitempty"`),
	}
}

// Edges of the Trip.
func (Trip) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("trips").
			Required().
			Unique().
			Field("user_id"),
		edge.From("driver", VehicleDriver.Type).
			Ref("trips").
			Unique().
			Field("driver_id"),
	}
}
