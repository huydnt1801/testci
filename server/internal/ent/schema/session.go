package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Session struct {
	ent.Schema
}

func (Session) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id").
			Optional(),
		field.Int("driver_id").
			Optional(),
		field.String("session_id"),
		field.Time("expire_in"),
		field.Bool("revoked").
			Default(false),
	}
}

func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("sessions").
			Unique().
			Field("user_id"),
		edge.From("driver", VehicleDriver.Type).
			Ref("sessions").
			Unique().
			Field("driver_id"),
	}
}

func (Session) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("session_id"),
	}
}
