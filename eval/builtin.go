package eval

//func newWindow(args ...object.Object) object.Object {
//
//	if len(args) != 3 {
//		return argumentAmountError(3, len(args))
//	}
//
//	width, ok := args[0].(*object.Integer)
//	if !ok {
//		return argumentTypeError(object.INTEGER_OBJ, args[0].Type(), 1)
//	}
//
//	height, ok := args[1].(*object.Integer)
//	if !ok {
//		return argumentTypeError(object.INTEGER_OBJ, args[1].Type(), 2)
//	}
//
//	title, ok := args[2].(*object.String)
//	if !ok {
//		return argumentTypeError(object.STR_OBJ, args[2].Type(), 3)
//	}
//
//	window := guilib.InitNewMainWindow(title.Value, height.Value, width.Value)
//
//	widget := guilib.AddCentralWidget(window)
//
//	return &object.Window{
//		Width:      width.Value,
//		Height:     height.Value,
//		Title:      title.Value,
//		MainWindow: window,
//		MainWidget: widget,
//		Closed:     make(chan bool),
//	}
//}

//func addButton(args ...object.Object) object.Object {
//	if len(args) != 3 {
//		return argumentAmountError(3, len(args))
//	}
//
//	window, ok := args[0].(*object.Window)
//
//	if !ok {
//		return argumentTypeError(object.WINDOW_OBJ, args[0].Type(), 1)
//	}
//	text, ok := args[1].(*object.String)
//
//	if !ok {
//		return argumentTypeError(object.STR_OBJ, args[1].Type(), 2)
//	}
//
//	function, ok := args[2].(*object.Function)
//	if !ok {
//		return argumentTypeError(object.FUNCTION_OBJ, args[2].Type(), 3)
//	}
//
//	button := guilib.AddButton(text.Value, window.MainWidget)
//
//	button.ConnectClicked(func(bool) {
//		Eval(function.Body, function.Env)
//	})
//
//	return NULL
//}
//
//func openWindow(args ...object.Object) object.Object {
//	if len(args) != 1 {
//		return argumentAmountError(1, len(args))
//	}
//
//	window, ok := args[0].(*object.Window)
//	if !ok {
//		return argumentTypeError(object.WINDOW_OBJ, args[0].Type(), 1)
//	}
//
//	go func(w *object.Window) {
//
//		window.MainWindow.Show()
//
//	}(window)
//
//	<-window.Closed
//
//	window.MainWindow.Close()
//
//	return NULL
//}
//
//func closeWindow(args ...object.Object) object.Object {
//	if len(args) != 1 {
//		return argumentAmountError(1, len(args))
//	}
//
//	window, ok := args[0].(*object.Window)
//	if !ok {
//		return argumentTypeError(object.WINDOW_OBJ, args[0].Type(), 1)
//	}
//
//	window.Closed <- true
//
//	return NULL
//}
//
//func closeButton(args ...object.Object) object.Object {
//	if len(args) != 1 {
//		return argumentAmountError(1, len(args))
//	}
//
//	window, ok := args[0].(*object.Window)
//	if !ok {
//		return argumentTypeError(object.WINDOW_OBJ, args[0].Type(), 1)
//	}
//
//	button := guilib.AddButton("Close", window.MainWidget)
//
//	button.ConnectClicked(func(bool) {
//		window.Closed <- true
//	})
//	return NULL
//}
//
//func addEdit(args ...object.Object) object.Object {
//	if len(args) != 2 {
//		return argumentAmountError(2, len(args))
//	}
//
//	window, ok := args[0].(*object.Window)
//	if !ok {
//		return argumentTypeError(object.WINDOW_OBJ, args[0].Type(), 1)
//	}
//	text, ok := args[1].(*object.String)
//	if !ok {
//		return argumentTypeError(object.STR_OBJ, args[1].Type(), 2)
//	}
//
//	edit := guilib.AddEdit(text.Value, window.MainWidget)
//
//	return &object.Edit{
//		Edit: edit,
//	}
//}
//
//func addLabel(args ...object.Object) object.Object {
//	if len(args) != 2 {
//		return argumentAmountError(2, len(args))
//	}
//
//	window, ok := args[0].(*object.Window)
//	if !ok {
//		return argumentTypeError(object.WINDOW_OBJ, args[0].Type(), 1)
//	}
//
//	text, ok := args[1].(*object.String)
//	if !ok {
//		return argumentTypeError(object.STR_OBJ, args[0].Type(), 2)
//	}
//
//	label := guilib.AddLabel(text.Value, window.MainWidget)
//
//	return &object.Label{
//		Label: label,
//	}
//}

//func write(args ...object.Object) object.Object {
//	if len(args) != 2 {
//		return argumentAmountError(2, len(args))
//	}
//
//	switch arg := args[0].(type) {
//
//	//	case *object.Label:
//	//		arg.Label.SetText(args[1].Inspect())
//
//	default:
//		return newError("Cant To Value of Typ %s", args[0].Type())
//	}
//
//	return NULL
//
//}

//func read(args ...object.Object) object.Object {
//	if len(args) != 1 {
//		return argumentAmountError(1, len(args))
//	}
//
//	switch arg := args[0].(type) {
//	case *object.Edit:
//		return &object.String{
//			Value: arg.Edit.Text(),
//		}
//
//	case *object.Label:
//		return &object.String{
//			Value: arg.Label.Text(),
//		}
//
//	default:
//		return newError("Cant Read From Value of Typ %s", args[0].Type())
//	}
//
//}

//func maping(args ...Object) Object {
//	if len(args) != 2 {
//		return argumentAmountError(1, len(args))
//	}
//
//	arrObj, ok := args[0].(*Array)
//	if !ok {
//		return argumentTypeError(ARR_OBJ, args[0].Type(), 1)
//	}
//
//	fn, ok := args[1].(*Function)
//
//	if !ok {
//
//		return argumentTypeError(FUNCTION_OBJ, args[1].Type(), 2)
//	}
//
//	arr := arrObj.Elements
//	elements := []Object{}
//
//	f := func(fn *Function, args []Object) Object {
//		extEnv := extendFunctionEnv(fn, args)
//		evaluated := Eval(fn.Body, extEnv)
//
//		return unwrapReturnValue(evaluated)
//	}
//
//	for _, obj := range arr {
//		elements = append(elements, f(fn, []Object{obj}))
//	}
//
//	return &Array{
//		Elements: elements,
//	}
//}
