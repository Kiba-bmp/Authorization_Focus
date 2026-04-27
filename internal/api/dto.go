package api

// JSON field names match the Android Retrofit DTOs (camelCase).

type RegisterRequest struct {
	UserName string `json:"userName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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
