package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type PcapWorkspace struct {
	ent.Schema
}

func (PcapWorkspace) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.UUID("user_id", uuid.UUID{}),
		field.String("file_name"),
		field.Enum("status").
			Values("LISTENING", "COMPLETED", "WAITING"),
	}
}

func (PcapWorkspace) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("pcapworkspaces").
			Field("user_id").
			Unique().
			Required(),
		edge.To("sessions", Session.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
