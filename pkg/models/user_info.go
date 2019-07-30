package models

type UserInfo struct {
	Name       string        `db:"name" json:"name"`
	Company string           `db:"company" json:"company"`
	Telephone string         `db:"telephone" json:"telephone"`
	Signed     int           `db:"signed" json:"signed"`
	Mark     string          `db:"mark" json:"mark"`
	Conference    int        `db:"conference" json:"conference"`
	Meeting  string          `db:"meeting" json:"meeting"`
}
