// Package org defines organization-level domain entities and relationships.
package org

// UserFunction associates a user with a specific function within the organization.
type UserFunction struct {
	UserID     int  `db:"user_id"`
	FunctionID int  `db:"function_id"`
	UnitID     int  `db:"unit_id"`
	RankID     *int `db:"rank_id"`
}
