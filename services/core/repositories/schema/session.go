package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Session struct {
	ent.Schema
}

func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique(),
		field.Time("sent_at").Default(time.Now()),
		field.String("client_host"),
		field.Uint32("server_port"),
		field.Uint8("ttl"),
		field.Int("packets_count"),
		field.String("client_user_agent").Optional(),
		field.String("protocol"),
		field.Float("average_response_time"),
		field.Int("requests_number"),
	}
}

func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("workspace", Workspace.Type).
			Ref("sessions").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("pcap_workspace", PcapWorkspace.Type).
			Ref("sessions"),
		edge.To("in_cart", CartSession.Type).
			Unique().
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("packets", Packet.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("session_filters", SessionFilter.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
