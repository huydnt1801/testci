package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type VehicleDriver struct {
	ent.Schema
}

func (VehicleDriver) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

func (VehicleDriver) Fields() []ent.Field {
	return []ent.Field{
		field.String("phone_number").
			Unique().
			MaxLen(15),
		field.String("full_name"),
		field.String("password"),
		field.Enum("license").
			Values("motor", "car"),
	}
}

func (VehicleDriver) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("trips", Trip.Type),
		edge.To("sessions", Session.Type),
	}
}
