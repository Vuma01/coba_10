package response

type SignupResponse struct {
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password" binding:"required"`
}
