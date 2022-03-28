package mysql

import (
	"fmt"
	"strings"
)

type SQLUpdates struct {
	assignments []string
	values      []interface{}
}

func (s *SQLUpdates) add(key string, value interface{}) {
	if s.assignments == nil {
		s.assignments = make([]string, 0, 1)
	}
	if s.values == nil {
		s.values = make([]interface{}, 0, 1)
	}
	s.assignments = append(s.assignments, fmt.Sprintf("%s = ?", key))
	s.values = append(s.values, value)
}

func (s SQLUpdates) Assignments() string {
	return strings.Join(s.assignments, ", ")
}

func (s SQLUpdates) Values() []interface{} {
	return s.values
}
