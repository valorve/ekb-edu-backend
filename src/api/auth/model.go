package auth

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type ChangePasswordInfo struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
