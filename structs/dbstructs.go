package structs

import "database/sql"

type DBRosterStruct struct {
	OwnerID              string         `db:"user_id"`
	UserID               string         `db:"friend_id"`
	Nickname             string         `db:"nickname"`
	Relation             sql.NullString `db:"state"`
	ContactState         sql.NullString `db:"contact_state"`
	ContactStateDate     sql.NullString `db:"contact_state_date"`
	ContactStatusMessage sql.NullString `db:"contact_status_message"`
	Regdate              sql.NullString `db:"regdate"`
	AvatarID             sql.NullString `db:"avatar_id"`
}
