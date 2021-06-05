package enki

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type statementBuilder struct {
	*primitiveBuilder
}

func NewStatement() *statementBuilder {
	return &statementBuilder{NewPrimitive()}
}

func (sb *statementBuilder) Line(s string, args ...interface{}) {
	if sb.Err() != nil {
		return
	}
	r := regexp.MustCompile(`@\d+`)
	result := r.ReplaceAllStringFunc(s, func(sub string) string {
		if sb.Err() != nil {
			return sub
		}
		d, err := strconv.Atoi(strings.TrimPrefix(sub, "@"))
		if err != nil {
			sb.err = err
			return sub
		}
		if len(args) < d {
			sb.err = fmt.Errorf("%s: found %s substitution parameter, but only %d arguments has been provided", s, sub, len(args))
			return sub
		}
		return fmt.Sprint(args[d-1])
	})
	_, err := sb.inner.WriteString(result)
	if err != nil {
		sb.err = err
		return
	}
	sb.NewLine()
}
