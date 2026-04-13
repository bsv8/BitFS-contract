package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ProcDomainEvents 对齐表 proc_domain_events。
type ProcDomainEvents struct {
	ent.Schema
}

func (ProcDomainEvents) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "proc_domain_events"}}
}

func (ProcDomainEvents) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at_unix"),
		field.String("command_id"),
		field.String("gateway_pubkey_hex"),
		field.String("event_name"),
		field.String("state_before"),
		field.String("state_after"),
		field.String("payload_json"),
	}
}

func (ProcDomainEvents) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("command_id", "id"),
		index.Fields("created_at_unix"),
		index.Fields("gateway_pubkey_hex", "id"),
	}
}
