package user_dtos

type CreateUserDTO struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,containsany=uppercase,containsany=lowercase,containsany=numeric,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}
