package psql

import (
	"fmt"
	"strings"
)

type SQLUpdateBuilder struct {
	assignment []string
	values     []interface{}
	query      string
}

func NewSQLUpdateBuilder(tableName string) *SQLUpdateBuilder {
	return &SQLUpdateBuilder{
		assignment: make([]string, 0, 1),
		values:     make([]interface{}, 0, 1),
		query:      fmt.Sprintf("UPDATE %s SET ", tableName),
	}
}

func (s *SQLUpdateBuilder) AddUpdateColumn(column string, value interface{}) {
	s.assignment = append(s.assignment, fmt.Sprintf("%s = $%d", column, len(s.values)+1))
	s.values = append(s.values, value)
}

func (s *SQLUpdateBuilder) AddWhere(column string, value interface{}) {
	s.query += strings.Join(s.assignment, ", ")
	s.query += fmt.Sprintf(" WHERE %s = $%d", column, len(s.values)+1)
	s.values = append(s.values, value)
}

func (s *SQLUpdateBuilder) AddReturning(columns ...string) {
	s.query += " RETURNING "

	for i, col := range columns {
		s.query += col

		if i != len(columns)-1 {
			s.query += ", "
		}
	}
}

func (s *SQLUpdateBuilder) GetQuery() string {
	return s.query
}

func (s *SQLUpdateBuilder) GetValues() []interface{} {
	return s.values
}
