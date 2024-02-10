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
	Id           int
	FullName     string
	Phone        *string
	Email        string
	PasswordHash string
	Role         RoleType
	Institution  *Institution
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
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Role         string `json:"role"`
}

type UserTokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
