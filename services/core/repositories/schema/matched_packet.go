package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type MatchedPacket struct {
	ent.Schema
}

func (MatchedPacket) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.UUID("session_filter_id", uuid.UUID{}),
		field.Int("packet_idx").Default(1),
	}
}

func (MatchedPacket) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("session_filter", SessionFilter.Type).
			Ref("matched_packets").
			Unique().
			Field("session_filter_id").
			Required(),
	}
}

func (MatchedPacket) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id", "packet_idx").Unique(),
	}
}
