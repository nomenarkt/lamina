package admin

type CreateUserRequest struct {
	CompanyID int    `json:"company_id" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Role      string `json:"role" binding:"required"` // enforce required now for consistency
}
