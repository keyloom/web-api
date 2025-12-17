package token_dtos

type PasswordGrantRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	ClientID string `form:"client_id" binding:"required"`
}
