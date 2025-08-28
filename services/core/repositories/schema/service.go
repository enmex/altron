package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Service struct {
	ent.Schema
}

func (Service) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.UUID("user_id", uuid.UUID{}),
		field.String("name"),
		field.String("link"),
		field.Uint32("port"),
		field.String("container_id").Optional().Nillable(),
	}
}

func (Service) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("services").
			Field("user_id").
			Unique().
			Required(),
		edge.To("plugins", Plugin.Type),
		edge.To("workspaces", Workspace.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("filters", Filter.Type),
		edge.To("sessions", Session.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Service) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "port").Unique(),
		index.Fields("user_id", "name").Unique(),
	}
}
