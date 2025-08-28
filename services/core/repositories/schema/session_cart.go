package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type CartSession struct {
	ent.Schema
}

func (CartSession) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.UUID("workspace_id", uuid.UUID{}),
		field.UUID("session_id", uuid.UUID{}),
	}
}

func (CartSession) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("session", Session.Type).
			Ref("in_cart").
			Field("session_id").
			Unique().
			Required(),
		edge.From("workspace", Workspace.Type).
			Ref("cart_sessions").
			Field("workspace_id").
			Unique().
			Required(),
	}
}
