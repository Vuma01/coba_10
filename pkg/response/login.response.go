package response

type LoginResponse struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password" binding:"required"`
}
