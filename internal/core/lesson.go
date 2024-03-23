package core

type LessonModel struct {
	Id          int
	Title       string
	ClassroomId int
	Content     *[]LessonContent
	Active      bool
}

type UpdateLessonModel struct {
	Id          int
	Title       *string
	ClassroomId *int
	Content     *[]LessonContent
	Active      *bool
}

//type LessonContentModel struct {
//	LessonContentItems []LessonContentItemModel
//}
//
//type LessonContentItemModel struct {
//	Id              string
//	Type            string
//	ExtraAttributes map[string]interface{}
//}

type Lesson struct {
	Id          int
	Title       string
	ClassroomId int
	Content     *[]LessonContent
	Active      bool
}

type UpdateLesson struct {
	Id          int
	Title       *string
	ClassroomId *int
	Content     *[]LessonContent
	Active      *bool
}

type LessonResponse struct {
	Id          int              `json:"id"`
	Title       string           `json:"title"`
	ClassroomId int              `json:"classroom_id"`
	Content     *[]LessonContent `json:"content"`
	Active      bool             `json:"active"`
}

type LessonContent struct {
	Id              string                 `json:"id"`
	Type            string                 `json:"type"`
	ExtraAttributes map[string]interface{} `json:"extra_attributes,omitempty"`
}

type CreateLessonRequest struct {
	Title  string `json:"title"`
	Active bool   `json:"active"`
}

type UpdateLessonRequest struct {
	// TODO: change current classroom id
	CurrentClassroomId *int             `json:"classroom_id,omitempty"`
	LessonId           *int             `json:"lesson_id,omitempty"`
	Title              *string          `json:"title,omitempty"`
	ClassroomId        *int             `json:"classroom_id,omitempty"`
	Content            *[]LessonContent `json:"content,omitempty"`
	Active             *bool            `json:"active,omitempty"`
}
