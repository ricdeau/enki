package enki

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Statement builder, that can build any statement.
type Statement interface {
	fmt.Stringer
	Block
	// Print prints statement to inner buffer.
	Print(s string, args ...interface{}) Statement
	// Line add line with args.
	Line(s string, args ...interface{}) Statement
	NewLine() Statement
}

type statement struct {
	inner strings.Builder
	err   error
}

// Stmt creates new statement builder.
func Stmt() *statement {
	return &statement{
		inner: strings.Builder{},
	}
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
		d, _ := strconv.Atoi(strings.TrimPrefix(sub, "@"))
		if len(args) < d {
			s.err = fmt.Errorf("%s: found %s substitution parameter, but only %d arguments has been provided", line, sub, len(args))
			return sub
		}
		return fmt.Sprint(args[d-1])
	})
	s.inner.WriteString(result)

	return s
}

// NewLine add new empty line.
func (s *statement) NewLine() Statement {
	_, err := s.inner.WriteString("\n")
	if err != nil {
		s.err = err
	}

	return s
}

func (s *statement) String() string {
	return s.inner.String()
}

func (s *statement) materialize() string {
	return s.String()
}
