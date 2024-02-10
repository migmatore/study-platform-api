package core

type LessonModel struct {
	Id          int
	Title       string
	ClassroomId int
	Active      bool
}

type Lesson struct {
	Id          int
	Title       string
	ClassroomId int
	Active      bool
}

type LessonResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	ClassroomId int    `json:"classroom_id"`
	Active      bool   `json:"active"`
}
