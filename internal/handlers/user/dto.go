package user

type Register struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Login struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
