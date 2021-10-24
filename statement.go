package enki

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Statement builder, that can build any statement.
type Statement interface {
	Primitive
	Block
	Print(s string, args ...interface{}) Statement
	Line(s string, args ...interface{}) Statement
}

type statement struct {
	*primitiveBuilder
}

// Stmt creates new statement builder.
func Stmt() *statement {
	return &statement{NewPrimitive()}
}

func Field(line string, args ...interface{}) Statement {
	return Stmt().Line(line, args...)
}

func (s *statement) Line(line string, args ...interface{}) Statement {
	s.Print(line, args...).NewLine()
	return s
}

func (s *statement) Print(line string, args ...interface{}) Statement {
	if s.err != nil {
		return s
	}

	r := regexp.MustCompile(`@\d+`)
	result := r.ReplaceAllStringFunc(line, func(sub string) string {
		if s.err != nil {
			return sub
		}
		d, err := strconv.Atoi(strings.TrimPrefix(sub, "@"))
		if err != nil {
			s.err = err
			return sub
		}
		if len(args) < d {
			s.err = fmt.Errorf("%s: found %s substitution parameter, but only %d arguments has been provided", line, sub, len(args))
			return sub
		}
		return fmt.Sprint(args[d-1])
	})
	_, err := s.inner.WriteString(result)
	if err != nil {
		s.err = err
		return s
	}

	return s
}

func (s *statement) materialize() string {
	return s.String()
}
