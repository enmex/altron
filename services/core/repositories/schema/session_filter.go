package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type SessionFilter struct {
	ent.Schema
}

func (SessionFilter) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.UUID("session_id", uuid.UUID{}),
		field.UUID("filter_id", uuid.UUID{}),
		field.Int("matches_count").Default(1),
	}
}

func (SessionFilter) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("session", Session.Type).
			Ref("session_filters").
			Unique().
			Field("session_id").
			Required(),
		edge.From("filter", Filter.Type).
			Ref("session_filter").
			Unique().
			Field("filter_id").
			Required(),
		edge.To("matched_packets", MatchedPacket.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (SessionFilter) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("session_id", "filter_id").Unique(),
	}
}
