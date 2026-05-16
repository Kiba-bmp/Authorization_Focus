package api

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	UserID      string `json:"userId"`
	Email       string `json:"email"`
	AccessToken string `json:"accessToken"`
	ExpiresAt   string `json:"expiresAt"`
}

type ErrorBody struct {
	Error string `json:"error"`
}
