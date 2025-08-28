package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type AnalyzerPayload struct {
	ent.Schema
}

func (AnalyzerPayload) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.UUID("workspace_id", uuid.UUID{}),
		field.String("value"),
		field.Int("number"),
		field.UUID("analyzer_component_id", uuid.UUID{}),
	}
}

func (AnalyzerPayload) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("workspace", Workspace.Type).
			Ref("analyzer_payloads").
			Unique().
			Field("workspace_id").
			Required(),
		edge.From("analyzer_component", AnalyzerComponent.Type).
			Ref("analyzer_payloads").
			Unique().
			Field("analyzer_component_id").
			Required(),
	}
}

func (AnalyzerPayload) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("workspace_id", "analyzer_component_id", "value").Unique(),
		index.Fields("workspace_id", "analyzer_component_id"),
	}
}
