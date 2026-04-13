package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ProcObservedGatewayStates 对齐表 proc_observed_gateway_states。
type ProcObservedGatewayStates struct {
	ent.Schema
}

func (ProcObservedGatewayStates) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "proc_observed_gateway_states"}}
}

func (ProcObservedGatewayStates) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at_unix"),
		field.String("gateway_pubkey_hex"),
		field.String("source_ref"),
		field.Int64("observed_at_unix"),
		field.String("event_name"),
		field.String("state_before"),
		field.String("state_after"),
		field.String("pause_reason"),
		field.Int64("pause_need_satoshi"),
		field.Int64("pause_have_satoshi"),
		field.String("last_error"),
		field.String("payload_json"),
	}
}

func (ProcObservedGatewayStates) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_at_unix"),
		index.Fields("event_name", "id"),
		index.Fields("gateway_pubkey_hex", "id"),
		index.Fields("source_ref", "observed_at_unix", "id"),
		index.Fields("state_after", "id"),
	}
}
