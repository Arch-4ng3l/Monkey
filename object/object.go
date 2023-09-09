package object

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/Arch-4ng3l/Monkey/ast"
)

const (
	INTEGER_OBJ = "INTEGER"
	FLOAT_OBJ   = "FLOAT"
	BOOLEAN_OBJ = "BOOLEAN"
	STR_OBJ     = "STRING"
	ARR_OBJ     = "ARRAY"
	TIME_OBJ    = "TIME"

	WINDOW_OBJ = "WINDOW"
	EDIT_OBJ   = "EDIT"
	LABEL_OBJ  = "LABEL"

	NULL         = "NULL"
	RETURN_OBJ   = "RETURN_OBJ"
	FUNCTION_OBJ = "FUNCTION_OBJ"
	ERROR_OBJ    = "ERROR"
	BUILTIN_OBJ  = "BUILTIN_FUNCTION"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Time struct {
	Time   time.Time
	Format string
}

func (t *Time) Type() ObjectType {
	return TIME_OBJ
}

func (t *Time) Inspect() string {
	return fmt.Sprintf("%d", t.Time.UnixMilli())
}

//type Label struct {
//	Label *widgets.QLabel
//}
//
//func (l *Label) Type() ObjectType {
//	return LABEL_OBJ
//}
//
//func (l *Label) Inspect() string {
//	return "Label"
//}
//
//type Edit struct {
//	Edit *widgets.QLineEdit
//}
//
//func (e *Edit) Type() ObjectType {
//	return EDIT_OBJ
//}
//
//func (e *Edit) Inspect() string {
//	return "Edit"
//}
//
//type Window struct {
//	Height     int
//	Width      int
//	Title      string
//	MainWindow *widgets.QMainWindow
//	MainWidget *widgets.QWidget
//	Closed     chan bool
//}
//
//func (w *Window) Type() ObjectType {
//	return WINDOW_OBJ
//}
//
//func (w *Window) Inspect() string {
//	return fmt.Sprintf("Title:%s\nHeight:%d\nWidth:%d", w.Title, w.Height, w.Width)
//}

type BuiltInFunction func(args ...Object) Object

type BuiltIn struct {
	Fn BuiltInFunction
}

func (bi *BuiltIn) Type() ObjectType {
	return BUILTIN_OBJ
}
func (bi *BuiltIn) Inspect() string {
	return "built in function"
}

type Function struct {
	Params []*ast.Ident
	Body   *ast.BlockStatement
	Env    *Env
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	var params = []string{}
	for _, p := range f.Params {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}
func (e *Error) Inspect() string {
	return "ERROR : " + e.Message
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_OBJ
}
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ARR_OBJ
}

func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STR_OBJ
}
func (s *String) Inspect() string {
	return s.Value
}

type Integer struct {
	Value int
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType {
	return FLOAT_OBJ
}

func (f *Float) Inspect() string {
	return fmt.Sprintf("%f", f.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL
}
func (n *Null) Inspect() string {
	return "null"
}
