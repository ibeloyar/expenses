package model

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Login        string `json:"login"`
	UserID       int    `json:"userID"`
	UserRoleID   int    `json:"userRoleID"`
}

type RegistrationData struct {
	Login    string `json:"login" validate:"required,min=2,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4,max=20"`
}
