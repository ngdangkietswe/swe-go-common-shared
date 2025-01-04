package domain

type ResetPassword struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
