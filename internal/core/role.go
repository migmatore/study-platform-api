package core

type RoleType string

const (
	AdminRole   RoleType = "admin"
	TeacherRole RoleType = "teacher"
	StudentRole RoleType = "student"
)

type RoleModel struct {
	Id   int
	Name string
}
