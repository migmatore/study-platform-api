package core

type UserModel struct {
	Id            int
	FullName      string
	Phone         *string
	Email         string
	PasswordHash  string
	RoleId        int
	InstitutionId *int
}

type User struct {
	FullName     string
	Phone        *string
	Email        string
	PasswordHash string
	Role         RoleType
	Institution  *InstitutionModel
}

type UserSigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignupRequest struct {
	Email           string   `json:"email"`
	Password        string   `json:"password"`
	FullName        string   `json:"full_name"`
	InstitutionName string   `json:"institution_name,omitempty"`
	Role            RoleType `json:"role"`
}

type UserAuthResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}
