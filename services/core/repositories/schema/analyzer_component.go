package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type AnalyzerComponent struct {
	ent.Schema
}

func (AnalyzerComponent) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.String("name").Unique(),
	}
}

func (AnalyzerComponent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("analyzer_payloads", AnalyzerPayload.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
