package entity

type LoginRequest struct {
	Email    string
	Password string
}

type RegisterRequest struct {
	Email    string
	Username string
	Role     string
	Password string
}

type AccessToken struct {
	AccessToken string
	ExpiresAt   int64
}
