// Package crew defines data structures for crew assignments and related operations.
package crew

import "time"

// Assignment represents a crew member assigned to a flight.
type Assignment struct {
	ID           int64     `db:"id"`            // Unique identifier of the assignment
	FlightID     int64     `db:"flight_id"`     // Associated flight's ID
	CrewID       int       `db:"crew_id"`       // ID of the crew member (linked to users table)
	CrewRole     string    `db:"crew_role"`     // Crew role (e.g., CDB, OPL, CCA)
	InFunction   bool      `db:"in_function"`   // Whether the crew is in function or MEP
	PickupTime   time.Time `db:"pickup_time"`   // Pickup time
	CheckinTime  time.Time `db:"checkin_time"`  // Check-in time
	CheckoutTime time.Time `db:"checkout_time"` // Checkout time
	CreatedAt    time.Time `db:"created_at"`    // Record creation time
}

// AssignmentDetail represents enriched crew assignment info with flight and crew metadata.
type AssignmentDetail struct {
	ID            int       `db:"id" json:"id"`                         // Assignment ID
	CrewID        int       `db:"crew_id" json:"crew_id"`               // Crew member ID
	CrewRole      string    `db:"crew_role" json:"crew_role"`           // Role on flight
	InFunction    bool      `db:"in_function" json:"in_function"`       // In function status
	PickupTime    time.Time `db:"pickup_time" json:"pickup_time"`       // Pickup timestamp
	CheckinTime   time.Time `db:"checkin_time" json:"checkin_time"`     // Check-in timestamp
	CheckoutTime  time.Time `db:"checkout_time" json:"checkout_time"`   // Checkout timestamp
	FlightNumber  string    `db:"flight_number" json:"flight_number"`   // Flight number
	DepartureCode string    `db:"departure_code" json:"departure_code"` // Departure airport code
	ArrivalCode   string    `db:"arrival_code" json:"arrival_code"`     // Arrival airport code
	CrewName      string    `db:"crew_name" json:"crew_name"`           // Full name of crew member
	CrewEmail     string    `db:"crew_email" json:"crew_email"`         // Crew email address
}
