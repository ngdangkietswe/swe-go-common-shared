package util

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

// Ref: https://entgo.io/docs/schema-edges/#quick-summary

// One2One is a helper function to create a one-to-one edge.
func One2One[T any](to string, toType T) ent.Edge {
	return edge.To(to, toType).Unique()
}

// One2OneInverse is a helper function to create a one-to-one inverse edge.
func One2OneInverse[T any](from string, fromType T, ref, field string) ent.Edge {
	return edge.From(from, fromType).
		Ref(ref).
		Unique().
		Required().
		Field(field)
}

// One2Many is a helper function to create a one-to-many edge.
func One2Many[T any](to string, toType T) ent.Edge {
	return edge.To(to, toType)
}

// One2ManyInverse is a helper function to create a one-to-many inverse edge.
func One2ManyInverse[T any](from string, fromType T, ref, field string) ent.Edge {
	return edge.From(from, fromType).
		Ref(ref).
		Unique().
		Field(field)
}

// One2ManyInverseRequired is a helper function to create a one-to-many inverse edge with required field.
func One2ManyInverseRequired[T any](from string, fromType T, ref, field string) ent.Edge {
	return edge.From(from, fromType).
		Ref(ref).
		Unique().
		Required().
		Field(field)
}

// Many2Many is a helper function to create a many-to-many edge.
func Many2Many[T any](to string, toType T) ent.Edge {
	return edge.To(to, toType)
}

// Many2ManyInverse is a helper function to create a many-to-many inverse edge.
func Many2ManyInverse[T any](from string, fromType T, ref string) ent.Edge {
	return edge.From(from, fromType).
		Ref(ref)
}
