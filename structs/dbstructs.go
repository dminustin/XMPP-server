package structs

import "database/sql"

type DBRosterStruct struct {
	OwnerID              string         `db:"owner_id"`
	UserID               string         `db:"user_id"`
	Nickname             string         `db:"nickname"`
	Relation             sql.NullString `db:"relation"`
	ContactState         sql.NullString `db:"contact_state"`
	ContactStateDate     sql.NullString `db:"contact_state_date"`
	ContactStatusMessage sql.NullString `db:"contact_status_message"`
}
