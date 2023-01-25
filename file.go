package enki

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"strings"

	"golang.org/x/tools/imports"
)

// Block builder that can be materialized
type Block interface {
	materialize() string
}

// File builder for files
type File interface {
	Statement
	// GoFmt - manage go-fmt before write file. By default, it is enabled.
	GoFmt(enabled bool) File
	// GoImports - formats and adjusts imports for the provided file. By default, it is disabled.
	GoImports(enabled bool) File
	// GeneratedBy - add special comment to identify generated file.
	GeneratedBy(tool string)
	// Package - add this file package name.
	Package(pkg string)
	// Import - add import with optional alias.
	Import(alias, path string)
	// Add - add a code block.
	Add(b Block)
	// Vars - add vars to `var` block.
	Vars(s ...Statement)
	// Consts - add consts to `const` block.
	Consts(s ...Statement)
	// Output - file as raw bytes.
	Output() ([]byte, error)
	// Write - write output or dest io.Writer.
	Write(dest io.Writer) error
	// Create - create new file or rewrite existing.
	Create(fileName string) error
}

var _ File = &file{}

type file struct {
	*statement
	goFmtEnabled     bool
	goImportsEnabled bool
	generatedComment string
	pkg              string
	imports          [][2]string
	blocks           []Block
	vars             []Statement
	consts           []Statement
}

// NewFile creates new file.
func NewFile() *file {
	return &file{
		goFmtEnabled: true,
		statement:    Stmt(),
	}
}

func (f *file) GoFmt(enabled bool) File {
	f.goFmtEnabled = enabled
	return f
}

func (f *file) GoImports(enabled bool) File {
	f.goImportsEnabled = enabled
	return f
}

func (f *file) GeneratedBy(tool string) {
	f.generatedComment = "// Code generated by " + tool + ". DO NOT EDIT."
}

func (f *file) Package(pkg string) {
	f.pkg = pkg
}

func (f *file) Import(alias, path string) {
	f.imports = append(f.imports, [2]string{alias, path})
}

func (f *file) Add(b Block) {
	f.blocks = append(f.blocks, b)
}

func (f *file) Vars(s ...Statement) {
	f.vars = append(f.vars, s...)
}

func (f *file) Consts(s ...Statement) {
	f.consts = append(f.consts, s...)
}

func (f *file) NewLine() Statement {
	stmt := Stmt()
	stmt.NewLine()
	f.blocks = append(f.blocks, stmt)

	return f
}

func (f *file) Line(s string, args ...interface{}) Statement {
	f.blocks = append(f.blocks, Stmt().Line(s, args...))
	return f
}

func (f *file) Output() (rawOut []byte, err error) {
	if f.err != nil {
		return nil, f.err
	}

	if f.generatedComment != "" {
		f.statement.Line(f.generatedComment)
	}

	if f.pkg != "" {
		f.statement.Line("package @1", f.pkg)
	}

	if len(f.imports) > 0 {
		f.statement.Line("import (")
		for _, s := range f.imports {
			f.statement.Line(`@1 "@2"`, s[0], strings.Trim(s[1], " "))
		}
		f.statement.Line(")")
	}

	f.statement.NewLine()

	if len(f.consts) > 0 {
		f.statement.Line("const (")
		for _, s := range f.consts {
			f.statement.Line(s.materialize())
		}
		f.statement.Line(")")
	}

	f.statement.NewLine()

	if len(f.vars) > 0 {
		f.statement.Line("var (")
		for _, s := range f.vars {
			f.statement.Line(s.materialize())
		}
		f.statement.Line(")")
	}

	f.statement.NewLine()

	buf := bytes.NewBufferString(f.String())

	for _, block := range f.blocks {
		buf.WriteString(block.materialize())
	}

	return buf.Bytes(), nil
}

func (f *file) Write(dest io.Writer) (err error) {
	output, err := f.Output()
	if err != nil {
		return fmt.Errorf("make file's output: %v", err)
	}

	if f.goFmtEnabled {
		output, err = format.Source(output)
		if err != nil {
			return fmt.Errorf("go fmt file: %v", err)
		}
	}

	_, err = dest.Write(output)
	if err != nil {
		return fmt.Errorf("write to dest: %v", err)
	}

	return nil
}

func (f *file) Create(fileName string) error {
	output, err := f.Output()
	if err != nil {
		return fmt.Errorf("make file's output: %v", err)
	}

	if f.goFmtEnabled {
		output, err = format.Source(output)
		if err != nil {
			return fmt.Errorf("go fmt file: %v", err)
		}
	}

	if f.goImportsEnabled {
		output, err = imports.Process(fileName, output, nil)
		if err != nil {
			return fmt.Errorf("process go imports: %v", err)
		}
	}

	out, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("create file: %v", err)
	}
	defer out.Close()

	if _, err = out.Write(output); err != nil {
		return fmt.Errorf("write: %v", err)
	}

	return nil
}
