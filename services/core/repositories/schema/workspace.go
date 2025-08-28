package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Workspace struct {
	ent.Schema
}

func (Workspace) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.String("name"),
		field.Enum("status").
			Values("LISTENING", "COMPLETED", "WAITING"),
	}
}

func (Workspace) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("sessions", Session.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("service", Service.Type).
			Ref("workspaces").
			Unique(),
		edge.To("analyzer_payloads", AnalyzerPayload.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("cart_sessions", CartSession.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
