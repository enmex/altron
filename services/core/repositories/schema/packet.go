package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Packet struct {
	ent.Schema
}

func (Packet) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.UUID("session_id", uuid.UUID{}),
		field.Time("sent_at"),
		field.Bool("is_request"),
		field.String("payload"),
	}
}

func (Packet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("session", Session.Type).
			Ref("packets").
			Unique().
			Field("session_id").
			Required(),
	}
}
