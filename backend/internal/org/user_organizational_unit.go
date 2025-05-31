// Package org defines organization-level domain entities and relationships.
package org

// UserOrganizationalUnit associates a user with an organizational unit.
type UserOrganizationalUnit struct {
	UserID int `db:"user_id"`
	UnitID int `db:"unit_id"`
}
