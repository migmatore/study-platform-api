package repository

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
	"github.com/migmatore/study-platform-api/internal/repository/psql"
	"github.com/migmatore/study-platform-api/pkg/logger"
	"github.com/migmatore/study-platform-api/pkg/utils"
)

type ClassroomRepo struct {
	logger logger.Logger
	pool   psql.AtomicPoolClient
}

func NewClassroomRepo(logger logger.Logger, pool psql.AtomicPoolClient) *ClassroomRepo {
	return &ClassroomRepo{logger: logger, pool: pool}
}

func (r ClassroomRepo) Create(ctx context.Context, classroom core.ClassroomModel) (core.ClassroomModel, error) {
	q := `INSERT INTO classrooms(title, description, teacher_id, max_students) VALUES($1, $2, $3, $4) 
			RETURNING id, title, description, teacher_id, max_students`

	newCLassroom := core.ClassroomModel{}

	if err := r.pool.QueryRow(
		ctx,
		q,
		classroom.Title,
		classroom.Description,
		classroom.TeacherId,
		classroom.MaxStudents,
	).Scan(
		&newCLassroom.Id,
		&newCLassroom.Title,
		&newCLassroom.Description,
		&newCLassroom.TeacherId,
		&newCLassroom.MaxStudents,
	); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return core.ClassroomModel{}, err
		}

		r.logger.Errorf("Query error. %v", err)
		return core.ClassroomModel{}, err
	}

	return newCLassroom, nil
}

func (r ClassroomRepo) Delete(ctx context.Context, id int) error {
	q := `DELETE FROM classrooms WHERE id = $1`

	if _, err := r.pool.Exec(ctx, q, id); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return err
		}

		r.logger.Errorf("Query error. %v", err)
		return err
	}

	return nil
}

func (r ClassroomRepo) ById(ctx context.Context, id int) (core.ClassroomModel, error) {
	q := `SELECT id, title, description, teacher_id, max_students FROM classrooms WHERE id = $1`

	var classroom core.ClassroomModel

	if err := r.pool.QueryRow(ctx, q, id).Scan(
		&classroom.Id,
		&classroom.Title,
		&classroom.Description,
		&classroom.TeacherId,
		&classroom.MaxStudents,
	); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return classroom, err
		}

		r.logger.Errorf("Query error. %v", err)
		return classroom, err
	}

	return classroom, nil
}

func (r ClassroomRepo) TeacherClassrooms(ctx context.Context, teacherId int) ([]core.ClassroomModel, error) {
	q := `SELECT id, title, description, teacher_id, max_students FROM classrooms WHERE teacher_id = $1`

	classrooms := make([]core.ClassroomModel, 0)

	rows, err := r.pool.Query(ctx, q, teacherId)
	if err != nil {
		r.logger.Errorf("Query error. %v", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		classroom := core.ClassroomModel{}

		err := rows.Scan(
			&classroom.Id,
			&classroom.Title,
			&classroom.Description,
			&classroom.TeacherId,
			&classroom.MaxStudents,
		)
		if err != nil {
			r.logger.Errorf("Query error. %v", err)
			return nil, err
		}

		classrooms = append(classrooms, classroom)
	}

	return classrooms, nil
}

func (r ClassroomRepo) StudentClassrooms(ctx context.Context, studentId int) ([]core.ClassroomModel, error) {
	q := `SELECT c.id, c.title, c.description, c.teacher_id, c.max_students FROM classroom_students 
    	JOIN public.classrooms c ON c.id = classroom_students.classroom_id WHERE student_id = $1`

	classrooms := make([]core.ClassroomModel, 0)

	rows, err := r.pool.Query(ctx, q, studentId)
	if err != nil {
		r.logger.Errorf("Query error. %v", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		classroom := core.ClassroomModel{}

		err := rows.Scan(
			&classroom.Id,
			&classroom.Title,
			&classroom.Description,
			&classroom.TeacherId,
			&classroom.MaxStudents,
		)
		if err != nil {
			r.logger.Errorf("Query error. %v", err)
			return nil, err
		}

		classrooms = append(classrooms, classroom)
	}

	return classrooms, nil
}

func (r ClassroomRepo) IsIn(ctx context.Context, classroomId, studentId int) (bool, error) {
	q := `SELECT EXISTS(SELECT * FROM classroom_students WHERE classroom_id = $1 AND student_id = $2)`

	var isIn bool

	if err := r.pool.QueryRow(ctx, q, classroomId, studentId).Scan(&isIn); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return isIn, err
		}

		r.logger.Errorf("Query error. %v", err)
		return isIn, err
	}

	return isIn, nil
}

func (r ClassroomRepo) Students(ctx context.Context, classroomId int) ([]core.UserModel, error) {
	q := `SELECT u.id, u.full_name, u.phone, u.email, u.password_hash, u.role_id, u.institution_id FROM classroom_students 
    	JOIN public.users u on u.id = classroom_students.student_id WHERE classroom_id = $1`

	users := make([]core.UserModel, 0)

	rows, err := r.pool.Query(ctx, q, classroomId)
	if err != nil {
		r.logger.Errorf("Query error. %v", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		user := core.UserModel{}

		err := rows.Scan(
			&user.Id,
			&user.FullName,
			&user.Phone,
			&user.Email,
			&user.PasswordHash,
			&user.RoleId,
			&user.InstitutionId,
		)
		if err != nil {
			r.logger.Errorf("Query error. %v", err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r ClassroomRepo) StudentsByClassroomsId(ctx context.Context, ids []int) ([]core.StudentModel, error) {
	q := `select u.id, u.full_name, u.phone, u.email, array_agg(c.id) classrooms from classroom_students 
    		join public.users u on u.id = classroom_students.student_id
			join public.classrooms c on c.id = classroom_students.classroom_id
			where classroom_id = any ($1) group by u.id, u.full_name, u.phone, u.email;`

	students := make([]core.StudentModel, 0)

	rows, err := r.pool.Query(ctx, q, ids)
	if err != nil {
		r.logger.Errorf("Query error. %v", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		student := core.StudentModel{}

		err := rows.Scan(
			&student.Id,
			&student.FullName,
			&student.Phone,
			&student.Email,
			&student.ClassroomsId,
		)
		if err != nil {
			r.logger.Errorf("Query error. %v", err)
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}
