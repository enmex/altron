package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Filter struct {
	ent.Schema
}

func (Filter) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("service_id", uuid.UUID{}).Optional(),
		field.String("name"),
		field.String("regex").Optional(),
		field.Uint8("ttl").Optional(),
		field.Int("total_packets").Optional(),
		field.Bool("is_blocking"),
		field.Bool("in_request"),
		field.Bool("in_response"),
		field.String("color"),
	}
}

func (Filter) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("filters").
			Field("user_id").
			Unique().
			Required(),
		edge.From("service", Service.Type).
			Ref("filters").
			Field("service_id").
			Unique(),
		edge.To("session_filter", SessionFilter.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Filter) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "name").Unique(),
	}
}
