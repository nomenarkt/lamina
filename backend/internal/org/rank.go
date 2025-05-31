// Package org defines organization-level domain entities and relationships.
package org

// Rank represents the crew rank within an organization.
type Rank struct {
	ID         int    `db:"id"`
	FunctionID int    `db:"function_id"`
	Name       string `db:"name"`
	Level      *int   `db:"level"`
}
