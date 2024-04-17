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
	WSToken      string `json:"ws_token"`
	RefreshToken string `json:"refresh_token"`
	Role         string `json:"role"`
}

type UserTokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UserAuthRequest struct {
	Token string `json:"token"`
}

type StudentModel struct {
	Id           int
	FullName     string
	Phone        *string
	Email        string
	ClassroomsId []int
}

type Student struct {
	Id           int
	FullName     string
	Phone        *string
	Email        string
	ClassroomsId []int
}

type StudentResponse struct {
	Id           int     `json:"id"`
	FullName     string  `json:"full_name"`
	Phone        *string `json:"phone,omitempty"`
	Email        string  `json:"email"`
	ClassroomsId []int   `json:"classrooms_id,omitempty"`
}

type ProfileResponse struct {
	FullName string  `json:"full_name"`
	Phone    *string `json:"phone,omitempty"`
	Email    string  `json:"email"`
}
