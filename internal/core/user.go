package core

type UserModel struct {
	Id            int
	FullName      string
	Phone         string
	Email         string
	PasswordHash  string
	RoleId        int
	InstitutionId int
}

type UserSigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignupRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	InstitutionName string `json:"institutionName,omitempty"`
	RoleId          int    `json:"role_id"`
}
