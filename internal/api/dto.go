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
	UserName    string `json:"userName"`
	Email       string `json:"email"`
	InviteCode  string `json:"inviteCode"`
	AccessToken string `json:"accessToken"`
	ExpiresAt   string `json:"expiresAt"`
}

type RegisterPendingResponse struct {
	VerificationRequired bool   `json:"verificationRequired"`
	Email                string `json:"email"`
	Message              string `json:"message"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type UserMeResponse struct {
	UserID       string  `json:"userId"`
	UserName     string  `json:"userName"`
	Email        string  `json:"email"`
	AvatarEmoji  *string `json:"avatarEmoji,omitempty"`
	InviteCode   string  `json:"inviteCode"`
	TodayMinutes int     `json:"todayMinutes"`
	WeekMinutes  int     `json:"weekMinutes"`
}

type UpdateProfileRequest struct {
	CurrentPassword string  `json:"currentPassword" binding:"required"`
	UserName        *string `json:"userName,omitempty"`
	Email           *string `json:"email,omitempty"`
	NewPassword     *string `json:"newPassword,omitempty"`
	AvatarEmoji     *string `json:"avatarEmoji,omitempty"`
}

type ErrorBody struct {
	Error string `json:"error"`
}
