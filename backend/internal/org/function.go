// Package org defines organization-level domain entities and relationships.
package org

// Function represents a functional unit or responsibility in an organization.
type Function struct {
	ID          int     `db:"id"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
}
