package model

type AuthToken struct {
	Token string `validate:"required"`
}

type UserAuthData struct {
	ID    string
	Name  string
	Email string
	Roles []string
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	RefreshToken     string `json:"-"`
	RefreshExpAt     int64  `json:"-"`
	RefreshTokenName string `json:"-"`
	AccessToken      string `json:"access_token"`
	AccessExpAt      int64  `json:"-"`
}

type RefreshRequest struct {
	RefreshToken string `json:"-" validate:"required"`
}
