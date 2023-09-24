package object

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Arch-4ng3l/Monkey/ast"
	"github.com/Arch-4ng3l/Monkey/code"
)

var NullVal = &Null{}

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

	NULL                  = "NULL"
	RETURN_OBJ            = "RETURN_OBJ"
	FUNCTION_OBJ          = "FUNCTION_OBJ"
	ERROR_OBJ             = "ERROR"
	BUILTIN_OBJ           = "BUILTIN_FUNCTION"
	COMPILED_FUNCTION_OBJ = "COMPILED_FUNCTION"
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

type CompiledFunction struct {
	Instructions code.Instructions
	NumLocals    int
	NumParams    int
}

func (cf *CompiledFunction) Type() ObjectType {
	return COMPILED_FUNCTION_OBJ
}

func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
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

func GetBuiltIntBuName(name string) *BuiltIn {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}

var Builtins = []struct {
	Name    string
	Builtin *BuiltIn
}{
	{"len", &BuiltIn{Fn: length}},
	{"print", &BuiltIn{Fn: print}},
}

func length(args ...Object) Object {
	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	switch arg := args[0].(type) {
	case *String:
		return &Integer{
			Value: len(arg.Value),
		}
	case *Array:
		return &Integer{
			Value: len(arg.Elements),
		}
	default:
		return newError("No Supported Datatype %s", arg.Type())
	}
}

func push(args ...Object) Object {
	if len(args) != 2 {
		argumentAmountError(2, len(args))
	}

	arr, ok := args[0].(*Array)
	if !ok {

		return argumentTypeError(ARR_OBJ, args[0].Type(), 1)

	}
	length := len(arr.Elements)

	newElements := make([]Object, length+1, length+1)
	copy(newElements, arr.Elements)

	newElements[length] = args[1]

	return &Array{
		Elements: newElements,
	}
}

func print(args ...Object) Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}
	return NullVal
}

func sort(args ...Object) Object {
	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}
	arrObj, ok := args[0].(*Array)
	if !ok {
		return argumentTypeError(ARR_OBJ, args[0].Type(), 1)
	}
	var elements []*Integer
	for i, obj := range arrObj.Elements {
		if obj, ok := obj.(*Integer); ok {
			elements = append(elements, obj)
			continue
		} else {
			return newError("Element on Position %d is Not an INTEGER", i)
		}

	}

	quickSort(elements, 0, len(elements)-1)

	var objArr []Object
	for _, e := range elements {
		objArr = append(objArr, e)
	}

	return &Array{
		Elements: objArr,
	}
}

func quickSort(arr []*Integer, low, high int) {
	if low < high {
		pivot := partition(arr, low, high)
		quickSort(arr, low, pivot-1)
		quickSort(arr, pivot+1, high)
	}
}

func partition(arr []*Integer, low, high int) int {
	pivot := arr[high].Value
	i := low - 1

	for j := low; j < high; j++ {
		if arr[j].Value < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

func typeof(args ...Object) Object {
	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	if args[0] == nil {
		return newError("Invalid Argument")
	}

	return &String{
		Value: string(args[0].Type()),
	}
}

func randInt(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	i, ok := args[0].(*Integer)
	if !ok {
		return argumentTypeError(INTEGER_OBJ, args[0].Type(), 1)
	}

	return &Integer{
		Value: rand.Int() % i.Value,
	}
}
func randFloat(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	f, ok := args[0].(*Float)
	if !ok {
		return argumentTypeError(FLOAT_OBJ, args[0].Type(), 1)
	}

	return &Float{
		Value: rand.Float64() * (f.Value),
	}

}

func randIntArray(args ...Object) Object {

	if len(args) != 2 && len(args) != 3 {
		return argumentAmountError(3, len(args))
	}

	m, ok := args[0].(*Integer)
	if !ok {
		return argumentTypeError(INTEGER_OBJ, args[0].Type(), 1)
	}

	l, ok := args[1].(*Integer)
	if !ok {
		return argumentTypeError(INTEGER_OBJ, args[1].Type(), 2)
	}

	dupe := false
	if len(args) == 3 && m.Value <= l.Value {
		if arg, ok := args[2].(*Boolean); ok {
			dupe = arg.Value
		}
	}

	objs := []Object{}
	for i := 0; i < l.Value; i++ {

		obj := &Integer{
			Value: rand.Int() % m.Value,
		}

		if dupe {
			ok = true
			for j := 0; j < i; j++ {
				if obj.Value == objs[j].(*Integer).Value {
					i--
					ok = false
				}
			}
			if ok {
				objs = append(objs, obj)
			}
		} else {
			objs = append(objs, obj)
		}
	}

	return &Array{
		Elements: objs,
	}
}

func printf(args ...Object) Object {
	if len(args) < 1 {
		return argumentAmountError(1, len(args))
	}

	str, ok := args[0].(*String)
	if !ok {
		return argumentTypeError(STR_OBJ, args[0].Type(), 1)
	}
	list := []any{}
	for _, arg := range args[1:] {
		switch arg := arg.(type) {
		case *Integer:
			list = append(list, arg.Value)
		case *Boolean:
			list = append(list, arg.Value)
		case *Float:
			list = append(list, arg.Value)
		case *String:
			list = append(list, arg.Value)
		default:
			return newError("Object of Type %s cant be Formated", arg.Type())
		}
	}
	formatedStr := formatString(str.Value, list...)

	fmt.Println(formatedStr)

	return NullVal
}

func toStr(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	switch arg := args[0].(type) {
	case *Integer:
		return &String{
			Value: fmt.Sprintf("%d", arg.Value),
		}
	case *Float:
		return &String{
			Value: fmt.Sprintf("%f", arg.Value),
		}
	default:
		return newError("Object of Type %s cant be converted to String", arg.Type())
	}

}

func toInt(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	switch arg := args[0].(type) {
	case *String:
		if i, err := strconv.ParseInt(arg.Value, 0, 64); err == nil {
			return &Integer{
				Value: int(i),
			}
		} else {
			return newError("String cant be converted to Integer")
		}
	default:
	}

	return newError("Object of Type %s cant be converted to Integer", args[0].Type())
}

func toFloat(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	switch arg := args[0].(type) {
	case *String:
		if f, err := strconv.ParseFloat(arg.Value, 64); err == nil {
			return &Float{
				Value: f,
			}
		} else {
			return newError("String cant be converted to Integer")
		}
	default:
		return newError("Object of Type %s cant be converted to Integer", args[0].Type())
	}

}

func formatString(str string, args ...any) string {
	return fmt.Sprintf(str, args...)
}

func argumentAmountError(num1, num2 int) *Error {
	return newError("Want %d Arguments got %d", num1, num2)
}

func argumentTypeError(type1 string, type2 ObjectType, num int) *Error {
	return newError("Argument %d has to be of Type %s got %s", num, type1, type2)
}

func newError(format string, a ...interface{}) *Error {

	return &Error{
		Message: fmt.Sprintf(format, a...),
	}
}
