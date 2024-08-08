package models

type Blabber struct {
	Id              int
	Username        string `db:"username"`
	RealName        string `db:"real_name"`
	BlabName        string `db:"blab_name"`
	CreatedDate     string `db:"created_at"`
	NumberListeners int
	NumberListening int
}
