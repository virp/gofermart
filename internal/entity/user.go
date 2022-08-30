package entity

type User struct {
	ID           string  `db:"id"`
	Login        string  `db:"login"`
	PasswordHash string  `db:"password_hash"`
	Balance      float64 `db:"balance"`
	Withdrawn    float64 `db:"withdrawn"`
}
