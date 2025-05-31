// Package org defines organization-level domain entities and relationships.
package org

import "time"

// OrganizationalUnit represents a business or functional unit within the organization.
type OrganizationalUnit struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Type      string    `db:"type"` // direction, department, service, section
	ParentID  *int      `db:"parent_id"`
	CreatedBy *int      `db:"created_by"`
	CreatedAt time.Time `db:"created_at"`
}
