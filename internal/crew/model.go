package crew

import "time"

type CrewAssignment struct {
	ID           int64     `db:"id"`
	FlightID     int64     `db:"flight_id"`
	CrewID       int       `db:"crew_id"`     // corresponds to users.company_id
	CrewRole     string    `db:"crew_role"`   // e.g., CDB, OPL, CCA, etc.
	InFunction   bool      `db:"in_function"` // true = working, false = MEP
	PickupTime   time.Time `db:"pickup_time"`
	CheckinTime  time.Time `db:"checkin_time"`
	CheckoutTime time.Time `db:"checkout_time"`
	CreatedAt    time.Time `db:"created_at"`
}

type CrewAssignmentDetail struct {
	ID            int       `db:"id" json:"id"`
	CrewID        int       `db:"crew_id" json:"crew_id"`
	CrewRole      string    `db:"crew_role" json:"crew_role"`
	InFunction    bool      `db:"in_function" json:"in_function"`
	PickupTime    time.Time `db:"pickup_time" json:"pickup_time"`
	CheckinTime   time.Time `db:"checkin_time" json:"checkin_time"`
	CheckoutTime  time.Time `db:"checkout_time" json:"checkout_time"`
	FlightNumber  string    `db:"flight_number" json:"flight_number"`
	DepartureCode string    `db:"departure_code" json:"departure_code"`
	ArrivalCode   string    `db:"arrival_code" json:"arrival_code"`
	CrewName      string    `db:"crew_name" json:"crew_name"`
	CrewEmail     string    `db:"crew_email" json:"crew_email"`
}
