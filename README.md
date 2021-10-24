# enki

[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/ricdeau/enki.svg)](https://pkg.go.dev/github.com/ricdeau/enki)
[![Build Status](https://app.travis-ci.com/ricdeau/enki.svg?branch=master)](https://travis-ci.com/github/ricdeau/enki)
[![codecov](https://codecov.io/gh/ricdeau/enki/branch/master/graph/badge.svg?token=R1SRKIVD5Z)](https://codecov.io/gh/ricdeau/enki)
[![Go Report Card](https://goreportcard.com/badge/github.com/ricdeau/enki)](https://goreportcard.com/report/github.com/ricdeau/enki)

`enki` is a package for easy go files generation

## Example

```go
func main() {
    f := NewFile()
    f.Package("enki")
    f.GeneratedBy("enki")
    f.Import("", "fmt")
    f.Line("// Number type redefined")
    
    f.Add(T("Number").Is(Int))
    f.NewLine()
    
    f.Add(T("NumStruct").Struct(
        Field("a, b " + Int),
    ))
    f.NewLine()
    
    f.Add(T("Num").Interface(
        Def("AsNumber").Returns(Str),
        Def("Sum").Params(Int).Returns(Int),
    ))
    f.NewLine()
    
    f.Add(M("AsNumber").Receiver("n Number").Body(
        Stmt().Line("return fmt.Sprint(n)"),
    ).Returns(Str))
    f.NewLine()
    
    f.Add(M("Sum").Receiver("n Number").Params("x int").Body(
        Stmt().Line("return int(n) + x"),
    ).Returns(Int))
    f.NewLine()
    
    f.Add(M("Sum").Receiver("ns NumStruct").Body(
        Stmt().Line("return ns.a + ns.b"),
    ).Returns(Int))
    f.NewLine()
    
    f.Add(F("sum").Params("a, b " + Float32).Body(
        Stmt().Line("return @1(a + b)", Float64),
    ).Returns("s " + Float64))

    _ = f.Write(os.Stdout)
}
```

Output:
```go
// Code generated by enki. DO NOT EDIT.
package enki

import (
	"fmt"
)

// Number type redefined
type Number int

type NumStruct struct {
	a, b int
}

type Num interface {
	AsNumber() string
	Sum(int) int
}

func (n Number) AsNumber() string {
	return fmt.Sprint(n)
}

func (n Number) Sum(x int) int {
	return int(n) + x
}

func (ns NumStruct) Sum() int {
	return ns.a + ns.b
}

func sum(a, b float32) (s float64) {
	return float64(a + b)
}
```