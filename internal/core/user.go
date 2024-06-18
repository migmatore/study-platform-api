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
	Id            int
	FullName      string
	Phone         *string
	Email         string
	PasswordHash  string
	Role          RoleType
	InstitutionId *int
}

type UserProfile struct {
	FullName string
	Phone    *string
	Email    string
}

type UserProfileModel struct {
	FullName string
	Phone    *string
	Email    string
}

type UpdateUserProfileModel struct {
	FullName *string
	Phone    *string
	Email    *string
	Password *string
}

type UpdateUserProfile struct {
	FullName *string
	Phone    *string
	Email    *string
	Password *string
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

type Teacher struct {
	Id       int
	FullName string
	Phone    *string
	Email    string
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

type UpdateProfileRequest struct {
	FullName *string `json:"full_name,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}

type CreateStudentRequest struct {
	FullName     string  `json:"full_name"`
	Phone        *string `json:"phone,omitempty"`
	Email        string  `json:"email"`
	Password     string  `json:"password"`
	ClassroomsId []int   `json:"classrooms_id"`
}

type TeacherResponse struct {
	Id       int     `json:"id"`
	FullName string  `json:"full_name"`
	Phone    *string `json:"phone,omitempty"`
	Email    string  `json:"email"`
}

type CreateTeacherRequest struct {
	FullName string  `json:"full_name"`
	Phone    *string `json:"phone,omitempty"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
}
