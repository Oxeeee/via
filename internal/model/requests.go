package model

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

type EmailVerifyRequest struct {
	Code string `json:"code" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

type RequestPasswordReset struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPassword struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RequestVerificationToken struct {
	Username string `json:"username" binding:"required,min=3"`
}
