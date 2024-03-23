package core

type ClassroomModel struct {
	Id          int
	Title       string
	Description *string
	TeacherId   int
	MaxStudents int
}

type Classroom struct {
	Id          int
	Title       string
	Description *string
	TeacherId   int
	MaxStudents int
}

type ClassroomResponse struct {
	Id          int     `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	TeacherId   int     `json:"teacher_id"`
	MaxStudents int     `json:"max_students"`
}

type CreateClassroomRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	MaxStudents int     `json:"max_students"`
}
