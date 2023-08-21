package eval

import (
	"fmt"
	"math/rand"

	"github.com/Arch-4ng3l/Monkey/guilib"
	"github.com/Arch-4ng3l/Monkey/object"
)

func length(args ...object.Object) object.Object {
	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{
			Value: len(arg.Value),
		}
	case *object.Array:
		return &object.Integer{
			Value: len(arg.Elements),
		}
	default:
		return newError("No Supported Datatype %s", arg.Type())
	}
}

func push(args ...object.Object) object.Object {
	if len(args) != 2 {
		argumentAmountError(2, len(args))
	}

	arr, ok := args[0].(*object.Array)
	if !ok {

		return argumentTypeError(object.ARR_OBJ, args[0].Type(), 1)

	}
	length := len(arr.Elements)

	newElements := make([]object.Object, length+1, length+1)
	copy(newElements, arr.Elements)

	newElements[length] = args[1]

	return &object.Array{
		Elements: newElements,
	}
}

func print(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}
	return NULL
}

func maping(args ...object.Object) object.Object {
	if len(args) != 2 {
		return argumentAmountError(1, len(args))
	}

	arrObj, ok := args[0].(*object.Array)
	if !ok {
		return argumentTypeError(object.ARR_OBJ, args[0].Type(), 1)
	}

	fn, ok := args[1].(*object.Function)

	if !ok {

		return argumentTypeError(object.FUNCTION_OBJ, args[1].Type(), 2)
	}

	arr := arrObj.Elements
	elements := []object.Object{}

	f := func(fn *object.Function, args []object.Object) object.Object {
		extEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extEnv)

		return unwrapReturnValue(evaluated)
	}

	for _, obj := range arr {
		elements = append(elements, f(fn, []object.Object{obj}))
	}

	return &object.Array{
		Elements: elements,
	}
}

func sort(args ...object.Object) object.Object {
	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}
	arrObj, ok := args[0].(*object.Array)
	if !ok {
		return argumentTypeError(object.ARR_OBJ, args[0].Type(), 1)
	}
	var elements []*object.Integer
	for i, obj := range arrObj.Elements {
		if obj, ok := obj.(*object.Integer); ok {
			elements = append(elements, obj)
			continue
		} else {
			return newError("Element on Position %d is Not an INTEGER", i)
		}

	}

	quickSort(elements, 0, len(elements)-1)

	var objArr []object.Object
	for _, e := range elements {
		objArr = append(objArr, e)
	}

	return &object.Array{
		Elements: objArr,
	}
}

func quickSort(arr []*object.Integer, low, high int) {
	if low < high {
		pivot := partition(arr, low, high)
		quickSort(arr, low, pivot-1)
		quickSort(arr, pivot+1, high)
	}
}

func partition(arr []*object.Integer, low, high int) int {
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

func typeof(args ...object.Object) object.Object {
	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	if args[0] == nil {
		return newError("Invalid Argument")
	}

	return &object.String{
		Value: string(args[0].Type()),
	}
}

func newWindow(args ...object.Object) object.Object {

	if len(args) != 3 {
		return argumentAmountError(3, len(args))
	}

	width, ok := args[0].(*object.Integer)
	if !ok {
		return argumentTypeError(object.INTEGER_OBJ, args[0].Type(), 1)
	}

	height, ok := args[1].(*object.Integer)
	if !ok {
		return argumentTypeError(object.INTEGER_OBJ, args[1].Type(), 2)
	}

	title, ok := args[2].(*object.String)
	if !ok {
		return argumentTypeError(object.STR_OBJ, args[2].Type(), 3)
	}

	window := guilib.InitNewMainWindow(title.Value, height.Value, width.Value)

	widget := guilib.AddCentralWidget(window)

	return &object.Window{
		Width:      width.Value,
		Height:     height.Value,
		Title:      title.Value,
		MainWindow: window,
		MainWidget: widget,
		Closed:     make(chan bool),
	}
}

func addButton(args ...object.Object) object.Object {
	if len(args) != 3 {
		return argumentAmountError(3, len(args))
	}

	window, ok := args[0].(*object.Window)

	if !ok {
		return argumentTypeError(object.WINDOW_OBJ, args[0].Type(), 1)
	}
	text, ok := args[1].(*object.String)

	if !ok {
		return argumentTypeError(object.STR_OBJ, args[1].Type(), 2)
	}

	function, ok := args[2].(*object.Function)
	if !ok {
		return argumentTypeError(object.FUNCTION_OBJ, args[2].Type(), 3)
	}

	button := guilib.AddButton(text.Value, window.MainWidget)

	button.ConnectClicked(func(bool) {
		Eval(function.Body, function.Env)
	})

	return NULL
}

func openWindow(args ...object.Object) object.Object {
	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	window, ok := args[0].(*object.Window)
	if !ok {
		return argumentTypeError(object.WINDOW_OBJ, args[0].Type(), 1)
	}

	go func(w *object.Window) {

		window.MainWindow.Show()

	}(window)

	<-window.Closed

	window.MainWindow.Close()

	return NULL
}

func closeWindow(args ...object.Object) object.Object {
	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	window, ok := args[0].(*object.Window)
	if !ok {
		return argumentTypeError(object.WINDOW_OBJ, args[0].Type(), 1)
	}

	window.Closed <- true

	return NULL
}

func closeButton(args ...object.Object) object.Object {
	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	window, ok := args[0].(*object.Window)
	if !ok {
		return argumentTypeError(object.WINDOW_OBJ, args[0].Type(), 1)
	}

	button := guilib.AddButton("Close", window.MainWidget)

	button.ConnectClicked(func(bool) {
		window.Closed <- true
	})
	return NULL
}

func addEdit(args ...object.Object) object.Object {
	if len(args) != 2 {
		return argumentAmountError(2, len(args))
	}

	window, ok := args[0].(*object.Window)
	if !ok {
		return argumentTypeError(object.WINDOW_OBJ, args[0].Type(), 1)
	}
	text, ok := args[1].(*object.String)
	if !ok {
		return argumentTypeError(object.STR_OBJ, args[1].Type(), 2)
	}

	edit := guilib.AddEdit(text.Value, window.MainWidget)

	return &object.Edit{
		Edit: edit,
	}
}

func addLabel(args ...object.Object) object.Object {
	if len(args) != 2 {
		return argumentAmountError(2, len(args))
	}

	window, ok := args[0].(*object.Window)
	if !ok {
		return argumentTypeError(object.WINDOW_OBJ, args[0].Type(), 1)
	}

	text, ok := args[1].(*object.String)
	if !ok {
		return argumentTypeError(object.STR_OBJ, args[0].Type(), 2)
	}

	label := guilib.AddLabel(text.Value, window.MainWidget)

	return &object.Label{
		Label: label,
	}
}

func write(args ...object.Object) object.Object {
	if len(args) != 2 {
		return argumentAmountError(2, len(args))
	}

	switch arg := args[0].(type) {

	case *object.Label:
		arg.Label.SetText(args[1].Inspect())

	default:
		return newError("Cant To Value of Typ %s", args[0].Type())
	}

	return NULL

}

func read(args ...object.Object) object.Object {
	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	switch arg := args[0].(type) {
	case *object.Edit:
		return &object.String{
			Value: arg.Edit.Text(),
		}

	case *object.Label:
		return &object.String{
			Value: arg.Label.Text(),
		}

	default:
		return newError("Cant Read From Value of Typ %s", args[0].Type())
	}

}

func randInt(args ...object.Object) object.Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	i, ok := args[0].(*object.Integer)
	if !ok {
		return argumentTypeError(object.INTEGER_OBJ, args[0].Type(), 1)
	}

	return &object.Integer{
		Value: rand.Int() % i.Value,
	}
}
func randFloat(args ...object.Object) object.Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	f, ok := args[0].(*object.Float)
	if !ok {
		return argumentTypeError(object.FLOAT_OBJ, args[0].Type(), 1)
	}

	return &object.Float{
		Value: rand.Float64() * (f.Value),
	}

}

func randIntArray(args ...object.Object) object.Object {

	if len(args) != 2 && len(args) != 3 {
		return argumentAmountError(3, len(args))
	}

	m, ok := args[0].(*object.Integer)
	if !ok {
		return argumentTypeError(object.INTEGER_OBJ, args[0].Type(), 1)
	}

	l, ok := args[1].(*object.Integer)
	if !ok {
		return argumentTypeError(object.INTEGER_OBJ, args[1].Type(), 2)
	}

	dupe := false
	if len(args) == 3 && m.Value <= l.Value {
		if arg, ok := args[2].(*object.Boolean); ok {
			dupe = arg.Value
		}
	}

	objs := []object.Object{}
	for i := 0; i < l.Value; i++ {

		obj := &object.Integer{
			Value: rand.Int() % m.Value,
		}

		if dupe {
			ok = true
			for j := 0; j < i; j++ {
				if obj.Value == objs[j].(*object.Integer).Value {
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

	return &object.Array{
		Elements: objs,
	}
}

func argumentAmountError(num1, num2 int) *object.Error {
	return newError("Want %d Arguments got %d", num1, num2)
}

func argumentTypeError(type1 string, type2 object.ObjectType, num int) *object.Error {
	return newError("Argument %d has to be of Type %s got %s", num, type1, type2)
}
