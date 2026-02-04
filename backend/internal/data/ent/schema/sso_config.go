package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// SSOConfig holds the schema definition for the SSOConfig entity.
type SSOConfig struct {
	ent.Schema
}

// Fields of the SSOConfig.
func (SSOConfig) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("enabled").
			Default(false),
		field.String("protocol").
			Default("saml2"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}
