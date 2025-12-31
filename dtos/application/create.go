package application_dtos

type CreateApplicationDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
